package transcons

import (
	"github.com/dm03514/test-engine/results"
	log "github.com/sirupsen/logrus"
)

type Conditions struct {
	Tcs []TransCon
}

func (c Conditions) Evaluate(result results.Result) results.Result {
	log.WithFields(log.Fields{
		"component": "conditions",
	}).Debug("Evaluate()")
	// loop over each con, decorating each one, if one fails return
	// Error Result
	for _, condition := range c.Tcs {
		result = condition.Evaluate(result)
		if result.Error() != nil {
			return result
		}
	}
	return result
}
