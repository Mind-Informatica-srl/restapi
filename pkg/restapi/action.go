package restapi

import (
	"net/http"
)

// Action provide the way to associate a HandlerFunc to the path, the method and authorization info
type Action struct {
	ActionFunc     http.HandlerFunc
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
}

func (action *Action) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action.ActionFunc(w, r)
}

func (action *Action) IsSkipAuth() bool {
	return action.SkipAuth
}

func (action *Action) GetPath() string {
	return action.Path
}

func (action *Action) GetMethod() string {
	return action.Method
}

func (action *Action) GetAuthorizations() []string {
	return action.Authorizations
}

// AbstractAction represents the set of instructions to be executed when the server receive a certain request, identified by path and method
// It determines if the client need to be authenticated and the set of authorizations needed to execute the set of instructions
type AbstractAction interface {
	// ServeHTTP execute the instructions of the action
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	// IsSkipAuth return true if the action don't need authentication to be executed
	IsSkipAuth() bool
	// GetPath return the path which the action respond at
	GetPath() string
	// GetMethod return the method associated to the action
	GetMethod() string
	// GetAuthorizations return the set of authorizations needed to execute the action
	GetAuthorizations() []string
}
