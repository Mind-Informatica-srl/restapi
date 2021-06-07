package actions

import (
	"net/http"

	"gorm.io/gorm"
)

type DBDeleteAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	ScopeDB        func(r *http.Request) (func(*gorm.DB) *gorm.DB, error)
	Delegate       DBDeleteDelegate
}

func (action *DBDeleteAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBDeleteAction) GetPath() string {
	return action.Path
}

func (action *DBDeleteAction) GetMethod() string {
	if action.Method != "" {
		return action.Method
	}
	return "DELETE"
}

func (action *DBDeleteAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBDeleteAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	id, err := action.Delegate.ExtractPK(r)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	element := action.Delegate.CreateObject()
	if err := action.Delegate.AssignPK(element, id); err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest, Data: element}
	}
	db := action.Delegate.ProvideDB()
	if action.ScopeDB != nil {
		if scope, err := action.ScopeDB(r); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		} else {
			db = db.Scopes(scope)
		}
	}
	if err := db.Delete(&element).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError, Data: element}
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// DBDeleteDelegate expose the functions needed by a DBDeleteAction
type DBDeleteDelegate interface {

	// ProvideDB provide the gorm pool
	ProvideDB() *gorm.DB

	// ExtractPK extract the model's primary key from the http request
	ExtractPK(r *http.Request) (interface{}, error)

	// CreateObject create the model object
	CreateObject() interface{}

	// AssignPK assign the primary key value to the model object properly
	AssignPK(element interface{}, pk interface{}) error
}
