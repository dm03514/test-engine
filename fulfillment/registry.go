package fulfillment

import (
	"fmt"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/results"
	"github.com/dm03514/test-engine/transcons"
)

type Fulfiller interface {
	Execute(rs results.Results) <-chan results.Result
	Name() string
}

type load func(
	f map[string]interface{},
	name string,
	a actions.Action,
	cs transcons.Conditions,
) (Fulfiller, error)

type Registry struct {
	m map[string]load
}

func (r Registry) Load(f map[string]interface{}, name string, a actions.Action, cs transcons.Conditions) (Fulfiller, error) {
	t := f["type"].(string)
	load, ok := r.m[t]
	if !ok {
		return nil, fmt.Errorf("Unable to parse fulfillment type %s", t)
	}
	return load(f, name, a, cs)
}

func NewRegistry() (Registry, error) {
	return Registry{
		m: map[string]load{
			// "noop": NewNoop,
		},
	}, nil
}
