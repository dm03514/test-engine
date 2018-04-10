// +build unit

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
	f := NewDefaultFactory()

	noop, err := fulfillment.NewNoop(
		nil,
		"dummystate",
		actions.Subprocess{
			CommandName: "echo",
			Args:        []string{"hello world!"},
		},
		transcons.Conditions{
			[]transcons.TransCon{
				transcons.IntEqual{
					UsingProperty: "returncode",
					ToEqual:       0,
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	e, err := f.New(
		Test{
			States: []State{
				noop,
			},
			Timeout: time.Duration(1 * time.Minute),
		},
	)

	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	err = e.Run(ctx)

	if err != nil {
		t.Error(err)
	}
}
