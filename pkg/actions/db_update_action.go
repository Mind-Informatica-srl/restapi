package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type DBUpdateAction DatabaseAction

func (action *DBUpdateAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBUpdateAction) GetPath() string {
	return action.Path
}

func (action *DBUpdateAction) GetMethod() string {
	return action.Method
}

func (action *DBUpdateAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBUpdateAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := action.Delegate.DBProvider()
	id, err := action.Delegate.PKExtractor(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
	if ok, err := action.Delegate.PKVerificator(element, id); !ok {
		pke := NewPKNotVerifiedError(element, id, err)
		http.Error(w, pke.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Save(element).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(element)
}
