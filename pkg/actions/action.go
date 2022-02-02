package actions

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ActionFunc func(w http.ResponseWriter, r *http.Request) *ActionError

// Action provide the way to associate a HandlerFunc to the path, the method and authorization info
type Action struct {
	ActionFunc     ActionFunc
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
}

func (action *Action) Serve(w http.ResponseWriter, r *http.Request) *ActionError {
	return action.ActionFunc(w, r)
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
	// Serve execute the instructions of the action
	Serve(w http.ResponseWriter, r *http.Request) *ActionError
	// IsSkipAuth return true if the action don't need authentication to be executed
	IsSkipAuth() bool
	// GetPath return the path which the action respond at
	GetPath() string
	// GetMethod return the method associated to the action
	GetMethod() string
	// GetAuthorizations return the set of authorizations needed to execute the action
	GetAuthorizations() []string
}

// ActionError is error type returned by an ActionFunc
type ActionError struct {
	Err    error
	Status int
	Data   interface{}
}

func (e ActionError) Error() string {
	return e.Err.Error()
}

func (e ActionError) Unwrap() error {
	return e.Err
}

// PrimaryKeyIntExtractor extract the int pk from the request's vars
func PrimaryKeyIntExtractor(r *http.Request, idName string) (int, error) {
	vars := mux.Vars(r)
	pk, err := strconv.Atoi(vars[idName])
	if err != nil {
		return 0, err
	}
	return pk, nil
}

var ErrorMissingIdName = fmt.Errorf("missing id name")

// PrimaryKeyStringExtractor extract the string pkj from the request's vars
func PrimaryKeyStringExtractor(r *http.Request, idName string) (string, error) {
	vars := mux.Vars(r)
	pk := vars[idName]
	if pk == "" {
		return "", ErrorMissingIdName
	}
	return pk, nil
}

// PrimaryKeySingleIntExtractorMap extract the unique int pk from the request's var and
// put it inside of map[string]interface{}
func PrimaryKeySingleIntExtractorMap(r *http.Request, idName string) (interface{}, error) {
	res := make(map[string]interface{})
	if value, err := PrimaryKeyIntExtractor(r, idName); err != nil {
		return nil, err
	} else {
		res[idName] = value
	}
	return res, nil
}

// PrimaryKeySingleStringExtractorMap extract the unique string pk from the request's var and
// put it inside of map[string]interface{}
func PrimaryKeySingleStringExtractorMap(r *http.Request, idName string) (interface{}, error) {
	res := make(map[string]interface{})
	if value, err := PrimaryKeyStringExtractor(r, idName); err != nil {
		return nil, err
	} else {
		res[idName] = value
	}
	return res, nil
}

// PrimaryKeyFullExtractorMap extract ALL the pk fields from the request's var and
// put them inside of map[string]interface{}
//
// intIdNames is the slice with the names of int pks and
// stringIdNames is the slice with the names of string pks
//
// Example:
// var multipleIdsExtractor = func(r *http.Request) (interface{}, error) {
// 		return actions.PrimaryKeyFullExtractorMap(r, []string{"AziendaID"}, []string{"Cognome", "Nome"})
// }
func PrimaryKeyFullExtractorMap(r *http.Request, intIdNames []string, stringIdNames []string) (interface{}, error) {
	res := make(map[string]interface{})
	for _, idName := range intIdNames {
		if value, err := PrimaryKeyIntExtractor(r, idName); err != nil {
			return nil, err
		} else {
			res[idName] = value
		}
	}
	for _, idName := range stringIdNames {
		if value, err := PrimaryKeyStringExtractor(r, idName); err != nil {
			return nil, err
		} else {
			res[idName] = value
		}
	}
	return res, nil
}
