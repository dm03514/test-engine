package engine

import (
	"context"
	"fmt"
	"github.com/dm03514/test-engine/ids"
	"github.com/dm03514/test-engine/results"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

// State is what this engine can execute!
type State interface {
	Execute(context.Context, results.Results) <-chan results.Result
	Name() string
}

// Test contains all states and global config and metadata
type Test struct {
	Name    string
	States  []State
	Timeout time.Duration
}

// Engine decorates a test and contains execution metadata
type Engine struct {
	Test

	currentState int
	rs           *results.Results

	recordStateDuration StateDurationRecorder
	recordTestDuration  TestDurationRecorder
}

// ExecuteState kicks off a state and keeps track of state execution states
func (e *Engine) ExecuteState(ctx context.Context) (State, <-chan results.Result) {
	s := e.States[e.currentState]
	log.WithFields(log.Fields{
		"component":           "Engine.ExecuteState()",
		"current_state_index": e.currentState,
		"state":               s.Name(),
		"execution_id":        ctx.Value(ids.Execution("execution_id")),
	}).Info("executing")
	c := s.Execute(ctx, *e.rs)
	e.currentState++
	return s, c
}

// IsLastState checks if there are anymore states
func (e *Engine) IsLastState() bool {
	log.Infof("IsLastState(), currState %d : len(states): %d", e.currentState, len(e.States))
	return e.currentState == len(e.States)
}

// Run kicks off a test execution
func (e Engine) Run(ctx context.Context) error {
	testID := uuid.NewV4()
	log.WithFields(log.Fields{
		"component":    "engine.Run()",
		"execution_id": testID.String(),
	}).Info("running_engine")
	testExecutionStart := time.Now()
	e.rs = results.New()

	ctx, cancel := context.WithTimeout(ctx, e.Timeout)
	defer cancel()

	ctx = context.WithValue(ctx, ids.Execution("execution_id"), testID)

engineloop:
	for {

		stateExecutionStart := time.Now()
		s, resultChan := e.ExecuteState(ctx)

	stateexecutionloop:
		for {
			select {
			case <-ctx.Done():
				return fmt.Errorf("context done()")
			case <-time.After(e.Timeout):
				return fmt.Errorf("timeout")
			case r, more := <-resultChan:
				log.WithFields(log.Fields{
					"component":    "engine.Run()",
					"execution_id": testID.String(),
					"more":         more,
				}).Debug("<-resultChan")

				if !more && e.IsLastState() {
					e.recordStateDuration(
						s.Name(),
						e.Test.Name,
						time.Now().Sub(stateExecutionStart),
						nil,
					)
					break engineloop
				}

				if !more {
					e.recordStateDuration(
						s.Name(),
						e.Test.Name,
						time.Now().Sub(stateExecutionStart),
						nil,
					)
					break stateexecutionloop
				}

				if r.Error() != nil {
					e.recordStateDuration(
						s.Name(),
						e.Test.Name,
						time.Now().Sub(stateExecutionStart),
						r.Error(),
					)
					e.recordTestDuration(
						e.Test.Name,
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
		e.Test.Name,
		time.Now().Sub(testExecutionStart),
		nil,
	)

	return nil
}
