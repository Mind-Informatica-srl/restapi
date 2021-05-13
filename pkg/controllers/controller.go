package controllers

import "github.com/Mind-Informatica-srl/restapi/pkg/actions"

type Controller struct {
	Path    string
	Actions []actions.AbstractAction
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
