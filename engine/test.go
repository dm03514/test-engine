package engine

import (
	"context"
	"fmt"
	"time"
)

type State interface {
	Execute() <-chan error
}

type Test struct {
	States  []State
	Timeout time.Duration
}

type Engine struct {
	Test
	currentState int
}

func (e *Engine) ExecuteState() <-chan error {
	s := e.States[e.currentState]
	c := s.Execute()
	e.currentState++
	return c
}

func (e *Engine) IsLastState() bool {
	return e.currentState == len(e.States)
}

func (e Engine) Run(ctx context.Context) error {
loop:
	for {
		s := e.ExecuteState()

		select {
		case <-ctx.Done():
			return fmt.Errorf("Context Done().")
		case <-time.After(e.Timeout):
			return fmt.Errorf("Timeout")
		case err := <-s:
			if err != nil {
				return err
			}

			// check if this is the last state
			if e.IsLastState() {
				break loop
			}
		}
	}

	return nil
}
