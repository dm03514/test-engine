package fulfillment

import (
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/transcons"
)

type Poller struct {
	actions.Action
	transcons.Conditions

	N        string
	Interval int
	Timeout  int
}

func (p Poller) Name() string {
	return p.N
}
