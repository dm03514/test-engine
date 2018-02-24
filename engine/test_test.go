package engine

import (
	"context"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/fulfillment"
	"github.com/dm03514/test-engine/transcons"
	"testing"
	"time"
)

func TestEngine_Run(t *testing.T) {
	test := Test{
		States: []State{
			fulfillment.NoopFulillment{
				Action: actions.Subprocess{
					CommandName: "echo",
					Args:        []string{"hello world!"},
				},
				Conditions: transcons.Conditions{
					[]transcons.TransCon{
						transcons.IntEqual{
							UsingProperty: "returncode",
							ToEqual:       0,
						},
					},
				},
			},
		},
		Timeout: time.Duration(1 * time.Minute),
	}

	e := Engine{
		Test: test,
	}
	ctx := context.Background()
	err := e.Run(ctx)

	if err != nil {
		t.Error(err)
	}
}
