package engine

import (
	"context"
	"fmt"
	"github.com/dm03514/test-engine/results"
	log "github.com/sirupsen/logrus"
	"time"
)

type State interface {
	Execute() <-chan results.Result
}

type Test struct {
	Name    string
	States  []State
	Timeout time.Duration
}

type Engine struct {
	Test
	currentState int
}

func (e *Engine) ExecuteState() <-chan results.Result {
	s := e.States[e.currentState]
	log.Infof("ExecuteState() %+v", s)
	c := s.Execute()
	e.currentState++
	return c
}

func (e *Engine) IsLastState() bool {
	log.Infof("IsLastState(), currState %d : len(states): %d", e.currentState, len(e.States))
	return e.currentState == len(e.States)
}

func (e Engine) Run(ctx context.Context) error {
	log.Infof("Run()")

engineloop:
	for {
		s := e.ExecuteState()

	stateexecutionloop:
		for {
			select {
			case <-ctx.Done():
				return fmt.Errorf("Context Done().")
			case <-time.After(e.Timeout):
				return fmt.Errorf("Timeout")
			case r, more := <-s:
				log.Infof("Read From state %+v. (more = %+v)", r, more)

				if !more && e.IsLastState() {
					break engineloop
				}

				if !more {
					break stateexecutionloop
				}

				if r.Error() != nil {
					return r.Error()
				}
			}
		}
	}

	return nil
}
