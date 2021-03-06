package fulfillment

import (
	"context"
	"fmt"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/results"
	"github.com/dm03514/test-engine/transcons"
	log "github.com/sirupsen/logrus"
)

type loaderFn func(f map[string]interface{}, name string, a actions.Action, cs transcons.Conditions) (Fulfiller, error)

const defaultType = "noop.Noop"

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
	t, ok := f["type"]

	if !ok {
		log.Warnf("Fulfillment type %q not found, falling back to default: %q", t, defaultType)
		t = defaultType
	}

	load, ok := r.m[t.(string)]
	if !ok {
		return nil, fmt.Errorf("Unable to find loader for fulfillment type %q", t)
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
