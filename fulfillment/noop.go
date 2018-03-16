package fulfillment

import (
	"context"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/ids"
	"github.com/dm03514/test-engine/results"
	"github.com/dm03514/test-engine/transcons"
	log "github.com/sirupsen/logrus"
)

type NoopFulfillment struct {
	a  actions.Action
	cs transcons.Conditions

	name string
}

func (n NoopFulfillment) Execute(ctx context.Context, rs results.Results) <-chan results.Result {
	c := make(chan results.Result)
	// execute the action in another go routine, run the conditions
	// against the result
	go func() {
		log.WithFields(log.Fields{
			"component":    "NoopFulfillment",
			"execution_id": ctx.Value(ids.Execution("execution_id")),
		}).Info("Execute()")

		r, err := n.a.Execute(ctx, rs)
		if err != nil {
			c <- results.ErrorResult{
				From: r,
				Err:  err,
			}
			close(c)
			return
		}

		r = n.cs.Evaluate(r)
		c <- r
		close(c)
	}()

	return c
}

func (n NoopFulfillment) Name() string {
	return n.name
}

func NewNoop(f map[string]interface{}, name string, a actions.Action, cs transcons.Conditions) (Fulfiller, error) {
	return NoopFulfillment{
		a:    a,
		cs:   cs,
		name: name,
	}, nil
}
