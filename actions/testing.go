package actions

import (
	"context"
	"github.com/dm03514/test-engine/results"
)

// DummyAction doesn't do anything!
type DummyAction struct{}

// Execute a noop action
func (da DummyAction) Execute(context.Context, results.Results) (results.Result, error) {
	return nil, nil
}
