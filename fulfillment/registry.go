package fulfillment

import (
	"context"
	"fmt"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/results"
	"github.com/dm03514/test-engine/transcons"
)

type loaderFn func(f map[string]interface{}, name string, a actions.Action, cs transcons.Conditions) (Fulfiller, error)

// Fulfiller executes and return results
type Fulfiller interface {
	Execute(ctx context.Context, rs results.Results) <-chan results.Result
	Name() string
}

// Registry contains fulfiller string identifiers mapped to loader functions
type Registry struct {
	m map[string]loaderFn
}

// Load parses a generic map into a fulfiller
func (r Registry) Load(f map[string]interface{}, name string, a actions.Action, cs transcons.Conditions) (Fulfiller, error) {
	t := f["type"].(string)
	load, ok := r.m[t]
	if !ok {
		return nil, fmt.Errorf("Unable to parse fulfillment type %s", t)
	}
	return load(f, name, a, cs)
}

// NewRegistry creates a usable registry
func NewRegistry() (Registry, error) {
	return Registry{
		m: map[string]loaderFn{
			"noop.Noop":   NewNoop,
			"poll.Poller": NewPoller,
		},
	}, nil
}
