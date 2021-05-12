package restapi

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

func (action *DBGetAllAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := action.DBProvider()
	db, err := QueryFilter(db, r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	list := action.ObjectCreator()
	paginationScope, page, pageSize := Paginate(r)
	var count int64
	//se è richiesta la paginazione
	if err := db.Scopes(paginationScope).Find(list).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if page > 0 && pageSize > 0 {
		// abbiamo la paginazione. Si calcola anche la count
		element := action.ObjectCreator()
		if err = db.Model(element).Count(&count).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res := Pager{
			TotalCount: count,
			Items: func() interface{} {
				return list
			},
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
