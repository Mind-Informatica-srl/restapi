package restapi

import (
	"net/http"
)

type DBDeleteAction DatabaseAction

func (action *DBDeleteAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBDeleteAction) GetPath() string {
	return action.Path
}

func (action *DBDeleteAction) GetMethod() string {
	return action.Method
}

func (action *DBDeleteAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBDeleteAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := action.DBProvider()
	id, err := action.PKExtractor(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	element := action.ObjectCreator()
	action.PKAssigner(element, id)
	if err := db.Delete(element).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
