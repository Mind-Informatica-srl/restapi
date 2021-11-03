package controllers

import "github.com/Mind-Informatica-srl/restapi/pkg/actions"

// PKUrlProvider provide the url part referred to the pk
type PKUrlProvider interface {
	PKUrl() string
}

// CreateModelController create a standard controller based on a delegate
// without auth
func CreateModelController(path string,
	delegate interface{},
) Controller {
	return CreateModelControllerWithAuth(path,
		delegate,
		[]string{},
		[]string{},
		[]string{},
		[]string{},
		[]string{},
	)
}

// CreateModelControllerWithAuth create a standard controller based on a delegate
// with auth
func CreateModelControllerWithAuth(path string,
	delegate interface{},
	getAllAuth []string,
	getOneAuth []string,
	insertAuth []string,
	updateAuth []string,
	deleteAuth []string,
) Controller {
	pkurl := "/{id}"
	if d, ok := delegate.(PKUrlProvider); ok {
		pkurl = d.PKUrl()
	}
	gad := delegate.(actions.DBGetAllDelegate)
	getAllAction := actions.DBGetAllAction{
		Method:         "GET",
		Path:           "",
		Delegate:       gad,
		Authorizations: getAllAuth,
	}

	god := delegate.(actions.DBGetOneDelegate)
	getOneAction := actions.DBGetOneAction{
		Method:         "GET",
		Path:           pkurl,
		Delegate:       god,
		Authorizations: getOneAuth,
	}

	id := delegate.(actions.DBInsertDelegate)
	insertAction := actions.DBInsertAction{
		Method:         "POST",
		Path:           "",
		Delegate:       id,
		Authorizations: insertAuth,
	}

	ud := delegate.(actions.DBUpdateDelegate)
	updateAction := actions.DBUpdateAction{
		Method:         "PUT",
		Path:           pkurl,
		Delegate:       ud,
		Authorizations: updateAuth,
	}

	dd := delegate.(actions.DBDeleteDelegate)
	deleteAction := actions.DBDeleteAction{
		Method:         "DELETE",
		Path:           pkurl,
		Delegate:       dd,
		Authorizations: deleteAuth,
	}

	return NewController(path, []actions.AbstractAction{
		&getAllAction,
		&getOneAction,
		&insertAction,
		&updateAction,
		&deleteAction,
	})
}
