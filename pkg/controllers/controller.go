package controllers

import "github.com/Mind-Informatica-srl/restapi/pkg/actions"

type Controller struct {
	Path    string
	Actions []actions.AbstractAction
}

// ModifyActions modifica le azioni e restituisce il controller
func (c Controller) ModifyActions(modifier func([]actions.AbstractAction) []actions.AbstractAction) Controller {
	c.Actions = modifier(c.Actions)
	return c
}

func (c *Controller) AddAction(action actions.AbstractAction) {
	c.Actions = append(c.Actions, action)
}

func NewController(path string, acts []actions.AbstractAction) Controller {
	return Controller{
		Path:    path,
		Actions: acts,
	}
}
