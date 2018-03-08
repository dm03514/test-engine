package actions

import (
	"fmt"
)

type loaderFn func(map[string]interface{}) (Action, error)

type Registry struct {
	m map[string]loaderFn
}

func (r Registry) Load(am map[string]interface{}) (Action, error) {
	t := am["type"].(string)
	loaderFn, ok := r.m[t]
	if !ok {
		return nil, fmt.Errorf("Unable to parse action type %s", t)
	}
	return loaderFn(am)
}

func NewRegistry() (Registry, error) {
	return Registry{
		m: map[string]loaderFn{
			"shell.Subprocess": NewSubprocessFromMap,
			"http.Http":        NewHttpFromMap,
		},
	}, nil
}
