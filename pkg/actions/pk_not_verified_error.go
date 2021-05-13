package actions

import "fmt"

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
