package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"
)

type DBUpdateAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	ScopeDB        func(r *http.Request) (func(*gorm.DB) *gorm.DB, error)
	Delegate       DBUpdateDelegate
}

func (action *DBUpdateAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBUpdateAction) GetPath() string {
	return action.Path
}

func (action *DBUpdateAction) GetMethod() string {
	if action.Method != "" {
		return action.Method
	}
	return "PUT"
}

func (action *DBUpdateAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBUpdateAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	id, err := action.Delegate.ExtractPK(r)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	element := action.Delegate.CreateObject()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	if err = json.Unmarshal(reqBody, element); err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}
	if ok, err := action.Delegate.VerifyPK(element, id); !ok {
		pke := NewPKNotVerifiedError(element, id, err)
		return &ActionError{Err: pke, Status: http.StatusBadRequest}
	}
	db := action.Delegate.ProvideDB()
	if action.ScopeDB != nil {
		if scope, err := action.ScopeDB(r); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		} else {
			db = db.Scopes(scope)
		}
	}
	if err := db.Save(element).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(element); err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}

	return nil
}

// DBDeleteDelegate expose the functions needed by a DBDeleteAction
type DBUpdateDelegate interface {

	// ProvideDB provide the gorm pool
	ProvideDB() *gorm.DB

	// CreateObject create the model object
	CreateObject() interface{}

	// VerifyPK check if the value of the primary key in the model object is equal to the passed primary key value
	VerifyPK(element interface{}, pk interface{}) (bool, error)

	// ExtractPK extract the model's primary key from the http request
	ExtractPK(r *http.Request) (interface{}, error)
}
