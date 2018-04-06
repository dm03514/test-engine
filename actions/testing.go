package actions

import (
	"context"
	"github.com/dm03514/test-engine/results"
)

type DummyAction struct{}

func (da DummyAction) Execute(context.Context, results.Results) (results.Result, error) {
	return nil, nil
}
