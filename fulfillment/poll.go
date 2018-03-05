package fulfillment

import (
	"fmt"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/results"
	"github.com/dm03514/test-engine/transcons"
	log "github.com/sirupsen/logrus"
	"time"
)

type Poller struct {
	a  actions.Action
	cs transcons.Conditions

	name string

	Interval time.Duration
	Timeout  time.Duration
}

func (p Poller) Name() string {
	return p.name
}

// Naively has a timeout, needs context, interval can execute longer than timeout
func (p Poller) Execute(rs results.Results) <-chan results.Result {
	c := make(chan results.Result)

	i := time.NewTicker(p.Interval)
	t := time.After(p.Timeout)

	go func() {
		log.Infof("Starting poller interval: %q, timeout: %q", p.Interval, p.Timeout)

	forloop:
		for {
			select {
			case <-t:
				c <- results.ErrorResult{
					Err: fmt.Errorf("Timeout %q exceeded", p.Timeout),
				}
				break forloop

			case <-i.C:
				log.Infof("Polling! @ %q", p.Interval)
				r, err := p.a.Execute(rs)
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
	}, nil
}
