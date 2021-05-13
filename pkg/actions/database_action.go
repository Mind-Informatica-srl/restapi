package actions

import (
	"github.com/Mind-Informatica-srl/restapi/pkg/delegate"
)

// DatabaseAction represent an action that do sometingh with the database
type DatabaseAction struct {
	Path           string
	Method         string
	SkipAuth       bool
	Authorizations []string
	Delegate       delegate.Delegate
}
