package actions

import (
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type DBGetAllAction DatabaseAction

func (action *DBGetAllAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBGetAllAction) GetPath() string {
	return action.Path
}

func (action *DBGetAllAction) GetMethod() string {
	return action.Method
}

func (action *DBGetAllAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBGetAllAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	db := action.Delegate.DBProvider()
	db, err := QueryFilter(db, r.URL)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	list := action.Delegate.ListCreator()
	paginationScope, page, pageSize := Paginate(r)
	var count int64
	//se è richiesta la paginazione
	if err := db.Scopes(paginationScope).Find(list).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}
	if page > 0 && pageSize > 0 {
		// abbiamo la paginazione. Si calcola anche la count
		element := action.Delegate.ObjectCreator()
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
