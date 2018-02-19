package actions

import "fmt"

type actionLoader func(map[string]interface{}) (Action, error)

type ActionRegistry struct {
	m map[string]actionLoader
}

func (ar ActionRegistry) Load(am map[string]interface{}) (Action, error) {
	t := am["type"].(string)
	loaderFn, ok := ar.m[t]
	if !ok {
		return nil, fmt.Errorf("Unable to parse action type %s", t)
	}
	return loaderFn(am)
}

func NewActionRegistry() (ActionRegistry, error) {
	return ActionRegistry{
		m: map[string]actionLoader{
			"shell.Subprocess": NewSubprocessFromMap,
		},
	}, nil
}
