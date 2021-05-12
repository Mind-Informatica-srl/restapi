package restapi

type Controller struct {
	Path    string
	Actions []AbstractAction
}

func (c *Controller) AddAction(action AbstractAction) {
	c.Actions = append(c.Actions, action)
}

func NewController(path string, actions []*Action) Controller {
	abstractActions := make([]AbstractAction, len(actions))
	for idx, v := range actions {
		abstractActions[idx] = v
	}
	return Controller{
		Path:    path,
		Actions: abstractActions,
	}
}

func NewAbstractController(path string, actions []AbstractAction) Controller {
	return Controller{
		Path:    path,
		Actions: actions,
	}
}
