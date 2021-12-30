package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"
)

type DBInsertAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	ScopeDB        func(r *http.Request) (func(*gorm.DB) *gorm.DB, error)
	Delegate       DBInsertDelegate
}

func (action *DBInsertAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBInsertAction) GetPath() string {
	return action.Path
}

func (action *DBInsertAction) GetMethod() string {
	if action.Method != "" {
		return action.Method
	}
	return "POST"
}

func (action *DBInsertAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBInsertAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	element, err := action.Delegate.CreateObject(r)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}

	if err = json.Unmarshal(reqBody, element); err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}

	db := action.Delegate.ProvideDB()
	if action.ScopeDB != nil {
		if scope, err := action.ScopeDB(r); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		} else {
			db = db.Scopes(scope)
		}
	}
	if err := db.Create(element).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(element); err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}

	return nil
}

// DBDeleteDelegate expose the functions needed by a DBDeleteAction
type DBInsertDelegate interface {

	// ProvideDB provide the gorm pool
	ProvideDB() *gorm.DB

	// CreateObject create the model object
	CreateObject(r *http.Request) (interface{}, error)
}
