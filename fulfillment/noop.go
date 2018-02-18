package fulfillment

import (
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/results"
	"github.com/dm03514/test-engine/transcons"
	log "github.com/sirupsen/logrus"
)

type NoopFulillment struct {
	actions.Action

	transcons.Conditions
}

func (n NoopFulillment) Execute() <-chan results.Result {
	log.Infof("NoopFulfillment()")
	c := make(chan results.Result)
	// execute the action in another go routine, run the conditions
	// against the result
	go func() {
		r, err := n.Action.Execute()
		if err != nil {
			c <- results.ErrorResult{
				From: r,
				Err:  err,
			}
			close(c)
		}

		r = n.Conditions.Evaluate(r)
		c <- r
		close(c)
	}()

	return c
}
