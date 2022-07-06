package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
)

type DBGetAllAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	ScopeDB        func(r *http.Request) (func(*gorm.DB) *gorm.DB, error)
	CustomizeList  func(*http.Request, interface{}) (interface{}, error)
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
	dbCount := action.Delegate.ProvideDB()
	if action.ScopeDB != nil {
		if scope, err := action.ScopeDB(r); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		} else {
			db = db.Scopes(scope)
			dbCount = dbCount.Scopes(scope)
		}
	}
	db, err := QueryFilter(db, r.URL)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	dbCount, err = QueryFilter(dbCount, r.URL)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	list, err := action.Delegate.CreateList(r)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	paginationScope, page, pageSize := Paginate(r)
	orderScope, _, _ := Ordinate(r)
	var count int64
	//se è richiesta la paginazione
	if err := db.Scopes(orderScope, paginationScope).Find(list).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}
	if action.CustomizeList != nil {
		if list, err = action.CustomizeList(r, list); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		}
	}
	accept := r.Header["Accept"]
	if page >= 0 && pageSize > 0 {
		// abbiamo la paginazione. Si calcola anche la count
		element, err := action.Delegate.CreateObject(r)
		if err != nil {
			return &ActionError{Err: err, Status: http.StatusBadRequest}
		}
		if err = dbCount.Model(element).Count(&count).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
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
	} else if len(accept) == 1 && accept[0] == "text/csv" {
		if err := gocsv.Marshal(list, w); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		}
		fileName := "export.csv"
		if fileName, err = action.Delegate.CSVFileName(r); err != nil {
			return &ActionError{Err: err, Status: http.StatusInternalServerError}
		}
		w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))

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
	CreateObject(r *http.Request) (interface{}, error)

	// CreateList create the model object list to be filled by the database interrogation
	CreateList(r *http.Request) (interface{}, error)

	// Return the csv file name
	CSVFileName(r *http.Request) (string, error)
}
