package fulfillment

import (
	"context"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/results"
	"testing"
	"time"
)

func TestPoller_Execute(t *testing.T) {
	interval, err := time.ParseDuration("100ms")
	if err != nil {
		t.Error(err)
	}
	timeout, err := time.ParseDuration("5m")
	if err != nil {
		t.Error(err)
	}

	stubResult := actions.StubResult{}
	p := Poller{
		a: actions.DummyAction{
			ReturnResult: stubResult,
		},
		Interval: interval,
		Timeout:  timeout,
	}

	resultChan := p.Execute(context.Background(), results.Results{})
	r := <-resultChan

	sr, ok := r.(actions.StubResult)

	if !ok {
		t.Errorf("Expected nil returned, received %+v", r)
	}

	if sr.Error() != nil {
		t.Error(sr.Error())
	}

}
