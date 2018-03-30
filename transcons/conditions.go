package transcons

import (
	"context"
	"github.com/dm03514/test-engine/ids"
	"github.com/dm03514/test-engine/results"
	log "github.com/sirupsen/logrus"
)

// Conditions is a collection of state transition conditions
type Conditions struct {
	Tcs []TransCon
}

// Evaluate 's a results against ALL conditions.  Returns after the first failure
func (c Conditions) Evaluate(ctx context.Context, result results.Result) results.Result {
	log.WithFields(log.Fields{
		"component":    "conditions",
		"execution_id": ctx.Value(ids.Execution("execution_id")),
	}).Debug("Evaluate()")
	// loop over each con, decorating each one, if one fails return
	// Error Result
	for _, condition := range c.Tcs {
		result = condition.Evaluate(ctx, result)
		if result.Error() != nil {
			return result
		}
	}
	return result
}
