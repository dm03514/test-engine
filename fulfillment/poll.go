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

type Poller struct {
	a  actions.Action
	cs transcons.Conditions

	t    string
	name string

	Interval time.Duration
	Timeout  time.Duration
}

func (p Poller) Name() string {
	return p.name
}

// Naively has a timeout, needs context, interval can execute longer than timeout
func (p Poller) Execute(ctx context.Context, rs results.Results) <-chan results.Result {
	c := make(chan results.Result)

	i := time.NewTicker(p.Interval)
	t := time.After(p.Timeout)

	go func() {
		log.WithFields(log.Fields{
			"component":    p.t,
			"interval":     p.Interval.String(),
			"timeout":      p.Timeout.String(),
			"execution_id": ctx.Value(ids.Execution("execution_id")),
		}).Info("starting_poller")

	forloop:
		for {
			select {
			case <-t:
				c <- results.ErrorResult{
					Err: fmt.Errorf("Timeout %q exceeded", p.Timeout),
				}
				break forloop

			case <-i.C:
				log.WithFields(log.Fields{
					"component":    p.t,
					"interval":     p.Interval.String(),
					"timeout":      p.Timeout.String(),
					"execution_id": ctx.Value(ids.Execution("execution_id")),
				}).Debug("polling!")
				r, err := p.a.Execute(ctx, rs)
				if err != nil {
					c <- results.ErrorResult{
						From: r,
						Err:  err,
					}
					break forloop
				}

				r = p.cs.Evaluate(r)

				// if there is NO error return
				// if not continue polling until success or timeout
				if r.Error() == nil {
					c <- r
					break forloop
				}
			}
		}

		close(c)
	}()

	return c
}

func NewPoller(f map[string]interface{}, name string, a actions.Action, cs transcons.Conditions) (Fulfiller, error) {
	i, err := time.ParseDuration(f["interval"].(string))
	if err != nil {
		return nil, err
	}

	t, err := time.ParseDuration(f["timeout"].(string))
	if err != nil {
		return nil, err
	}

	return Poller{
		a:        a,
		cs:       cs,
		name:     name,
		Interval: i,
		Timeout:  t,
		t:        f["type"].(string),
	}, nil
}
