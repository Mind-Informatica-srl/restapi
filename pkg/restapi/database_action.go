package restapi

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// PKNotVerifiedError is the error raised when the primary key of the element doesn't correspond with the one passed in the url
type PKNotVerifiedError struct {
	element interface{}
	pk      interface{}
}

func NewPKNotVerifiedError(element interface{}, pk interface{}) PKNotVerifiedError {
	return PKNotVerifiedError{
		element: element,
		pk:      pk,
	}
}

func (a PKNotVerifiedError) Error() string {
	return fmt.Sprintf("the primary key of the element: %v doesn't correspond with the one passed in the url: %v", a.element, a.pk)
}

// CreateObjectFunc represent the function used by DatabaseAction to instantiate its specific model
type CreateObjectFunc func() interface{}

// DBFunc represent the function used by DatabaseAction to retrieve the gorm.DB pointer
type DBFunc func() *gorm.DB

// ExtractPKFunc represent the function used by DatabaseAction to extract the primary key of the managed model from the http request
type ExtractPKFunc func(r *http.Request) (interface{}, error)

// VerifyPKFunc represent the function used by DatabaseAction to verify the primary key value of the managed model.
// It return error if the primary key value don't corresponde.
type VerifyPKFunc func(element interface{}, pk interface{}) *PKNotVerifiedError

// AssignPKFunc represent the function used by DatabaseAction to assign the primary key value to the managed model.
type AssignPKFunc func(element interface{}, pk interface{})

// DatabaseAction represent an action that do sometingh with the database
type DatabaseAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	ObjectCreator  CreateObjectFunc
	DBProvider     DBFunc
	PKExtractor    ExtractPKFunc
	PKVerificator  VerifyPKFunc
	PKAssigner     AssignPKFunc
}
