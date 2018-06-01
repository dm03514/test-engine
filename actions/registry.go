package actions

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

const defaultType = "noop.Noop"

type loaderFn func(map[string]interface{}) (Action, error)

// Registry maps action identifiers to parser functions
type Registry struct {
	m map[string]loaderFn
}

// Load an action from a generic map
func (r Registry) Load(am map[string]interface{}) (Action, error) {
	t, ok := am["type"]

	if !ok {
		log.Warnf("Action type %q not found, falling back to default: %q", t, defaultType)
		t = defaultType
	}

	loaderFn, ok := r.m[t.(string)]
	if !ok {
		return nil, fmt.Errorf("Unable to parse action type %s", t)
	}
	return loaderFn(am)
}

// NewRegistry initializes loadable actions
func NewRegistry() (Registry, error) {
	return Registry{
		m: map[string]loaderFn{
			"shell.Subprocess": NewSubprocessFromMap,
			"http.Http":        NewHTTPFromMap,
			"noop.Noop":        NewNoopFromMap,
		},
	}, nil
}
