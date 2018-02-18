package fulfillment

import (
	"github.com/dm03514/test-engine/actions"
	log "github.com/sirupsen/logrus"
)

type NoopFulillment struct {
	actions.Action
}

func (n NoopFulillment) Execute() <-chan actions.Result {
	log.Infof("NoopFulfillment()")
	r := make(chan actions.Result)
	// execute the action in another go routine, run the conditions
	// against the result
	go func() {
		close(r)
		/*
			r, err := n.Action.Execute()
			if err != nil {
				r <- err
			}

			err = n.Conditions.Evaluate(r)
			if err != nil {
				return err
			}

			close
		*/
	}()

	return r
}
