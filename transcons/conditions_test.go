package transcons

import (
	"context"
	"github.com/dm03514/test-engine/actions"
	"testing"
)

func TestConditions_Evaluate_TwoConditions(t *testing.T) {
	cs := Conditions{
		[]TransCon{
			IntEqual{
				UsingProperty: "returncode",
				ToEqual:       0,
			},
			StringEqual{
				UsingProperty: "output",
				ToEqual:       "hello!",
			},
		},
	}
	r := cs.Evaluate(context.Background(), actions.SubprocessResult{
		Output:     []byte("hello!"),
		Returncode: 0,
	})
	if r.Error() != nil {
		t.Errorf("Expected no error received; %+v", r.Error())
	}
}
