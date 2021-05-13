package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type DBInsertAction DatabaseAction

func (action *DBInsertAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBInsertAction) GetPath() string {
	return action.Path
}

func (action *DBInsertAction) GetMethod() string {
	return action.Method
}

func (action *DBInsertAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBInsertAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	element := action.Delegate.ObjectCreator()

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(reqBody, element); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db := action.Delegate.DBProvider()
	if err := db.Create(element).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(element)
}
