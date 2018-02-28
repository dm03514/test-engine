package engine

import (
	"context"
	"fmt"
	"github.com/dm03514/test-engine/results"
	log "github.com/sirupsen/logrus"
	"time"
)

type State interface {
	Execute(results.Results) <-chan results.Result
	Name() string
}

type Test struct {
	Name    string
	States  []State
	Timeout time.Duration
}

type Engine struct {
	Test

	currentState int
	rs           *results.Results

	recordStateDuration DurationRecorder
	recordTestDuration  DurationRecorder
}

func (e *Engine) ExecuteState() (State, <-chan results.Result) {
	s := e.States[e.currentState]
	log.Infof("ExecuteState() %+v", s)
	c := s.Execute(*e.rs)
	e.currentState++
	return s, c
}

func (e *Engine) IsLastState() bool {
	log.Infof("IsLastState(), currState %d : len(states): %d", e.currentState, len(e.States))
	return e.currentState == len(e.States)
}

func (e Engine) Run(ctx context.Context) error {
	log.Infof("Run()")
	testExecutionStart := time.Now()
	e.rs = results.New()

engineloop:
	for {

		stateExecutionStart := time.Now()
		s, resultChan := e.ExecuteState()

	stateexecutionloop:
		for {
			select {
			case <-ctx.Done():
				return fmt.Errorf("Context Done().")
			case <-time.After(e.Timeout):
				return fmt.Errorf("Timeout")
			case r, more := <-resultChan:
				log.Infof("Read From state %+v. (more = %+v)", r, more)

				if !more && e.IsLastState() {
					break engineloop
				}

				if !more {
					e.recordStateDuration(
						time.Now().Sub(stateExecutionStart),
						nil,
					)
					break stateexecutionloop
				}

				if r.Error() != nil {
					e.recordStateDuration(
						time.Now().Sub(stateExecutionStart),
						r.Error(),
					)
					e.recordTestDuration(
						time.Now().Sub(testExecutionStart),
						r.Error(),
					)
					return r.Error()
				}

				e.rs.Add(
					results.NamedResult{
						Name:   s.Name(),
						Result: r,
					},
				)
			}
		}
	}

	e.recordTestDuration(
		time.Now().Sub(testExecutionStart),
		nil,
	)

	return nil
}
