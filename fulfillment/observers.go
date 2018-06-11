package fulfillment

import (
	"context"
	"fmt"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/ids"
	"github.com/dm03514/test-engine/results"
	"github.com/dm03514/test-engine/transcons"
	log "github.com/sirupsen/logrus"
	"time"
)

// EventsObserver can observe any type of events
type EventsObserver struct {
	Timeout time.Duration

	cs transcons.Conditions

	name string
	t    string
}

// Name identifies this Event Observer
func (e EventsObserver) Name() string {
	return e.name
}

// Execute subscribes to an events stream
func (e EventsObserver) Execute(ctx context.Context, rs results.Results) <-chan results.Result {
	outChan := make(chan results.Result)

	go func() {
		timeoutChan := time.After(e.Timeout)

		log.WithFields(log.Fields{
			"component":    e.t,
			"timeout":      e.Timeout.String(),
			"execution_id": ctx.Value(ids.Execution("execution_id")),
		}).Info("starting_events_observer")

		// TODO RECIEVE OVER AN OBSERVABLE SPECIFIED
		// TODO EVERY RECEIVE SHOULD INVOKE THE TRANSITION CONDITIONS AGAINST OBSERVED EVENT
		// TODO RESULT, should contain the most recent AND the total currently seen

	forloop:
		for {
			select {
			case <-timeoutChan:
				outChan <- results.ErrorResult{
					Err: fmt.Errorf("Timeout %q exceeded", e.Timeout),
				}
				break forloop
			}
		}

		close(outChan)
	}()

	return outChan
}

// NewEventsObserver constructs a new one
func NewEventsObserver(f map[string]interface{}, name string, a actions.Action, cs transcons.Conditions) (Fulfiller, error) {

	t, err := time.ParseDuration(f["timeout"].(string))
	if err != nil {
		return nil, err
	}

	return EventsObserver{
		cs:      cs,
		name:    name,
		Timeout: t,
		t:       f["type"].(string),
	}, nil
}
