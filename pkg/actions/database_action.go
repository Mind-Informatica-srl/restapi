package actions

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Mind-Informatica-srl/restapi/pkg/delegate"
)

// PKNotVerifiedError is the error raised when the primary key of the element doesn't correspond with the one passed in the url
type PKNotVerifiedError struct {
	element interface{}
	pk      interface{}
	err     error
}

func NewPKNotVerifiedError(element interface{}, pk interface{}, err error) PKNotVerifiedError {
	return PKNotVerifiedError{
		element: element,
		pk:      pk,
		err:     err,
	}
}

func (e PKNotVerifiedError) Error() string {
	return fmt.Sprintf("the primary key of the element: %v doesn't correspond with the one passed in the url: %v", e.element, e.pk)
}

func (e PKNotVerifiedError) Unwrap() error {
	return e.err
}

// DatabaseAction represent an action that do sometingh with the database
type DatabaseAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	Delegate       delegate.Delegate
}

// PrimaryKeyIntExtractor permette di recuperare l'id dalla request
// oltre alla request si passa il nome del parametro che identifica l'id
func PrimaryKeyIntExtractor(r *http.Request, idName string) (int, error) {
	vars := mux.Vars(r)
	pk, err := strconv.Atoi(vars[idName])
	if err != nil {
		return 0, err
	}
	return pk, nil
}

var ErrorMissingIdName = fmt.Errorf("missing id name")

// PrimaryKeyIntExtractor permette di recuperare l'id dalla request
// oltre alla request si passa il nome del parametro che identifica l'id
func PrimaryKeyStringExtractor(r *http.Request, idName string) (string, error) {
	vars := mux.Vars(r)
	pk := vars[idName]
	if pk == "" {
		return "", ErrorMissingIdName
	}
	return pk, nil
}
