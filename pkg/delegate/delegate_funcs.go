package delegate

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CreateObjectFunc represent the function used by DatabaseAction to instantiate its specific model
type CreateObjectFunc func() interface{}

// DBFunc represent the function used by DatabaseAction to retrieve the gorm.DB pointer
type DBFunc func() *gorm.DB

// ExtractPKFunc represent the function used by DatabaseAction to extract the primary key of the managed model from the http request
type ExtractPKFunc func(r *http.Request) (interface{}, error)

// VerifyPKFunc represent the function used by DatabaseAction to verify the primary key value of the managed model.
// It return error if the primary key value don't corresponde.
type VerifyPKFunc func(element interface{}, pk interface{}) (bool, error)

// AssignPKFunc represent the function used by DatabaseAction to assign the primary key value to the managed model.
type AssignPKFunc func(element interface{}, pk interface{})

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

// PrimaryKeyIntExtractor extract the string pkj from the request's vars
func PrimaryKeyStringExtractor(r *http.Request, idName string) (string, error) {
	vars := mux.Vars(r)
	pk := vars[idName]
	if pk == "" {
		return "", ErrorMissingIdName
	}
	return pk, nil
}
