package transcons

import (
	"context"
	"github.com/dm03514/test-engine/results"
)

type TransCon interface {
	Evaluate(ctx context.Context, result results.Result) results.Result
}
