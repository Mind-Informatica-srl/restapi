package actions

import (
	"encoding/json"
	"net/http"
)

type DBGetOneAction DatabaseAction

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

func (action *DBGetOneAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := action.Delegate.DBProvider()
	id, err := action.Delegate.PKExtractor(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	element := action.Delegate.ObjectCreator()
	if err := db.First(element, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(element); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
