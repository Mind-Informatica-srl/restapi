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

func (action *DBInsertAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	element := action.Delegate.ObjectCreator()

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}

	if err = json.Unmarshal(reqBody, element); err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}

	db := action.Delegate.DBProvider()
	if err := db.Create(element).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(element)

	return nil
}
