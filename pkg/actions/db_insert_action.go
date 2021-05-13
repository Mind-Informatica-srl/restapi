package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"
)

type DBInsertAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	Delegate       DBInsertDelegate
}

func (action *DBInsertAction) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *DBInsertAction) GetPath() string {
	return action.Path
}

func (action *DBInsertAction) GetMethod() string {
	if action.Method != "" {
		return action.Method
	}
	return "POST"
}

func (action *DBInsertAction) GetAuthorizations() []string {
	return action.Authorizations
}

func (action *DBInsertAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	element := action.Delegate.CreateObject()

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}

	if err = json.Unmarshal(reqBody, element); err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}

	db := action.Delegate.ProvideDB()
	if err := db.Create(element).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(element)

	return nil
}

// DBDeleteDelegate expose the functions needed by a DBDeleteAction
type DBInsertDelegate interface {

	// ProvideDB provide the gorm pool
	ProvideDB() *gorm.DB

	// CreateObject create the model object
	CreateObject() interface{}
}
