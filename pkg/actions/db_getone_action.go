package actions

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type DBGetOneAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	Delegate       DBGetOneDelegate
}

func (action *DBGetOneAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBGetOneAction) GetPath() string {
	return action.Path
}

func (action *DBGetOneAction) GetMethod() string {
	return action.Method
}

func (action *DBGetOneAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBGetOneAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	db := action.Delegate.ProvideDB()
	id, err := action.Delegate.ExtractPK(r)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	element := action.Delegate.CreateObject()
	if err := db.First(element, id).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(element); err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}
	return nil
}

// DBGetOneDelegate expose the functions needed by a DBGetOneAction
type DBGetOneDelegate interface {

	// ProvideDB provide the gorm pool
	ProvideDB() *gorm.DB

	// CreateObject create the model object
	CreateObject() interface{}

	// ExtractPK extract the model's primary key from the http request
	ExtractPK(r *http.Request) (interface{}, error)
}
