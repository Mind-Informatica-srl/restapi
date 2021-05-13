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

func (action *DBGetOneAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	db := action.Delegate.DBProvider()
	id, err := action.Delegate.PKExtractor(r)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	element := action.Delegate.ObjectCreator()
	if err := db.First(element, id).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(element); err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}
	return nil
}
