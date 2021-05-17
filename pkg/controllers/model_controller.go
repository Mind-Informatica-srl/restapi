package controllers

import "github.com/Mind-Informatica-srl/restapi/pkg/actions"

// CreateModelController create a standard controller based on a delegate
func CreateModelController(path string, delegate interface{}) Controller {
	gad := delegate.(actions.DBGetAllDelegate)
	getAllAction := actions.DBGetAllAction{
		Method:   "GET",
		Path:     "",
		Delegate: gad,
	}

	god := delegate.(actions.DBGetOneDelegate)
	getOneAction := actions.DBGetOneAction{
		Method:   "GET",
		Path:     "/{id}",
		Delegate: god,
	}

	id := delegate.(actions.DBInsertDelegate)
	insertAction := actions.DBInsertAction{
		Method:   "POST",
		Path:     "",
		Delegate: id,
	}

	ud := delegate.(actions.DBUpdateDelegate)
	updateAction := actions.DBUpdateAction{
		Method:   "PUT",
		Path:     "/{id}",
		Delegate: ud,
	}

	dd := delegate.(actions.DBDeleteDelegate)
	deleteAction := actions.DBDeleteAction{
		Method:   "DELETE",
		Path:     "/{id}",
		Delegate: dd,
	}

	return NewController(path, []actions.AbstractAction{
		&getAllAction,
		&getOneAction,
		&insertAction,
		&updateAction,
		&deleteAction,
	})
}
