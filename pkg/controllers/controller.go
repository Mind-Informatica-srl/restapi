package controllers

import "github.com/Mind-Informatica-srl/restapi/pkg/actions"

type Controller struct {
	Path    string
	Actions []actions.AbstractAction
}

func (c *Controller) AddAction(action actions.AbstractAction) {
	c.Actions = append(c.Actions, action)
}

func NewController(path string, acts []*actions.Action) Controller {
	abstractActions := make([]actions.AbstractAction, len(acts))
	for idx, v := range acts {
		abstractActions[idx] = v
	}
	return Controller{
		Path:    path,
		Actions: abstractActions,
	}
}

func NewAbstractController(path string, acts []actions.AbstractAction) Controller {
	return Controller{
		Path:    path,
		Actions: acts,
	}
}
