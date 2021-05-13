package actions

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

func (action *DBDeleteAction) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	db := action.Delegate.DBProvider()
	id, err := action.Delegate.PKExtractor(r)
	if err != nil {
		return &ActionError{Err: err, Status: http.StatusBadRequest}
	}
	element := action.Delegate.ObjectCreator()
	action.Delegate.PKAssigner(element, id)
	if err := db.Delete(element).Error; err != nil {
		return &ActionError{Err: err, Status: http.StatusInternalServerError}
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
