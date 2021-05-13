package actions

import (
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type DBGetAllAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	ScopeDB        func(db *gorm.DB, r *http.Request) (func(*gorm.DB) *gorm.DB, error)
	Delegate       DBGetAllDelegate
}

func (action *DBGetAllAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBGetAllAction) GetPath() string {
	return action.Path
}

func (action *DBGetAllAction) GetMethod() string {
	if action.Method != "" {
		return action.Method
	}
	return "GET"
}

func (action *DBGetAllAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBGetAllAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	db := action.Delegate.ProvideDB()
	if action.ScopeDB != nil {
		if scope, err := action.ScopeDB(db, r); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		} else {
			db = db.Scopes(scope)
		}
	}
	db, err := QueryFilter(db, r.URL)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	list := action.Delegate.CreateList()
	paginationScope, page, pageSize := Paginate(r)
	var count int64
	//se è richiesta la paginazione
	if err := db.Scopes(paginationScope).Find(list).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}
	if page > 0 && pageSize > 0 {
		// abbiamo la paginazione. Si calcola anche la count
		element := action.Delegate.CreateObject()
		if err = db.Model(element).Count(&count).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		}
		res := Pager{
			TotalCount: count,
			Items: func() interface{} {
				return list
			},
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		}
	} else if err := json.NewEncoder(w).Encode(list); err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}
	return nil
}

// DBGetAllDelegate expose the functions needed by a DBDGetAllAction
type DBGetAllDelegate interface {

	// ProvideDB provide the gorm pool
	ProvideDB() *gorm.DB

	// CreateObject create the model object
	CreateObject() interface{}

	//CreateList create the model object list to be filled by the database interrogation
	CreateList() interface{}
}
