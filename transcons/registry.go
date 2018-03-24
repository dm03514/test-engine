package transcons

import (
	"fmt"
)

type loaderFn func(map[string]interface{}) (TransCon, error)

// Registry contains a map of string identifiers to initializer functions
type Registry struct {
	m map[string]loaderFn
}

// Load a registry based on a type
func (r Registry) Load(tcm map[string]interface{}) (TransCon, error) {
	t := tcm["type"].(string)
	loaderFn, ok := r.m[t]
	if !ok {
		return nil, fmt.Errorf("Unable to parse transaction condition type type %s", t)
	}
	return loaderFn(tcm)
}

// NewRegistry creates a registry of all available assertions
func NewRegistry() (Registry, error) {
	return Registry{
		m: map[string]loaderFn{
			"assertions.IntEqual":    NewIntEqualFromMap,
			"assertions.StringEqual": NewStringEqualFromMap,
			"assertions.Subprocess":  NewSubprocessFromMap,
		},
	}, nil
}
