package transcons

import (
	"context"
	"github.com/dm03514/test-engine/results"
)

// TransCon interface for evaluating state transition conditions
type TransCon interface {
	Evaluate(ctx context.Context, result results.Result) results.Result
}
