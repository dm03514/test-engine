package transcons

import (
	"github.com/dm03514/test-engine/results"
)

type TransCon interface {
	Evaluate(result results.Result) results.Result
}
