package actions

import (
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type DBGetOneAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	ScopeDB        func(r *http.Request) (func(*gorm.DB) *gorm.DB, error)
	Delegate       DBGetOneDelegate
}

func (action *DBGetOneAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBGetOneAction) GetPath() string {
	return action.Path
}

func (action *DBGetOneAction) GetMethod() string {
	if action.Method != "" {
		return action.Method
	}
	return "GET"
}

func (action *DBGetOneAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBGetOneAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	ids, err := action.Delegate.ExtractPK(r)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	element, err := action.Delegate.CreateObject(r)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	db := action.Delegate.ProvideDB()
	if action.ScopeDB != nil {
		if scope, err := action.ScopeDB(r); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		} else {
			db = db.Scopes(scope)
		}
	}
	var rnf bool
	if err := db.First(element, ids).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		rnf = true
	} else if err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}
	if !rnf {
		if err := json.NewEncoder(w).Encode(element); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		}
	}
	return nil
}

// DBGetOneDelegate expose the functions needed by a DBGetOneAction
type DBGetOneDelegate interface {

	// ProvideDB provide the gorm pool
	ProvideDB() *gorm.DB

	// CreateObject create the model object
	CreateObject(r *http.Request) (interface{}, error)

	// ExtractPK extract the model's primary key from the http request
	ExtractPK(r *http.Request) (interface{}, error)
}
