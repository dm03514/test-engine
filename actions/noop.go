package actions

import (
	"context"
	"github.com/dm03514/test-engine/results"
	"github.com/mitchellh/mapstructure"
)

// Noop action doesn't do anything!
type Noop struct {
	Type string
}

// Execute doesn't do anything and doesn't return anything useful!
func (n Noop) Execute(ctx context.Context, rs results.Results) (results.Result, error) {
	return nil, nil
}

// NewNoopFromMap creates a noop action
func NewNoopFromMap(m map[string]interface{}) (Action, error) {
	var n Noop
	err := mapstructure.Decode(m, &n)
	return n, err
}
