package actions

import (
	"context"
	"github.com/dm03514/test-engine/results"
)

// Action allows execution
type Action interface {
	Execute(context.Context, results.Results) (results.Result, error)
}
