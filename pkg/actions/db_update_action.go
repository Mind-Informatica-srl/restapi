package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"
)

type DBUpdateAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	Delegate       DBUpdateDelegate
}

func (action *DBUpdateAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBUpdateAction) GetPath() string {
	return action.Path
}

func (action *DBUpdateAction) GetMethod() string {
	if action.Method != "" {
		return action.Method
	}
	return "PUT"
}

func (action *DBUpdateAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBUpdateAction) Serve(w http.ResponseWriter, r *http.Request) {
	db := action.Delegate.ProvideDB()
	id, err := action.Delegate.ExtractPK(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	element := action.Delegate.CreateObject()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(reqBody, element); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ok, err := action.Delegate.VerifyPK(element, id); !ok {
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

// DBDeleteDelegate expose the functions needed by a DBDeleteAction
type DBUpdateDelegate interface {

	// ProvideDB provide the gorm pool
	ProvideDB() *gorm.DB

	// CreateObject create the model object
	CreateObject() interface{}

	// VerifyPK check if the value of the primary key in the model object is equal to the passed primary key value
	VerifyPK(element interface{}, pk interface{}) (bool, error)

	// ExtractPK extract the model's primary key from the http request
	ExtractPK(r *http.Request) (interface{}, error)
}
