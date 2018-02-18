package fulfillment

import "github.com/dm03514/test-engine/actions"

type NoopFulillment struct {
	actions.Action
}

func (n NoopFulillment) Execute() <-chan error {
	r := make(chan error)
	return r
}
