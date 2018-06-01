package observables

import "fmt"

type loaderFn func(map[string]interface{}) (Observable, error)

// Registry contains all observables that can be used in tests
type Registry struct {
	m map[string]loaderFn
}

// Load an observable from a map
func (r Registry) Load(am map[string]interface{}) (Observable, error) {
	t := am["type"].(string)
	loaderFn, ok := r.m[t]
	if !ok {
		return nil, fmt.Errorf("Unable to parse action type %s", t)
	}
	return loaderFn(am)
}

// NewRegistry initializes loadable actions
func NewRegistry() (Registry, error) {
	return Registry{
		m: map[string]loaderFn{
			"observables.HTTP": NewHTTPObservableFromMap,
		},
	}, nil
}
