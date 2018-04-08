package fulfillment

import (
	"context"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/results"
	"testing"
)

func TestNoopFulfillment_Execute(t *testing.T) {
	f := NoopFulfillment{
		a: actions.DummyAction{},
	}
	outChan := f.Execute(context.Background(), results.Results{})
	r := <-outChan
	if r != nil {
		t.Errorf("Expected nil return received %+v", r)
	}
}
