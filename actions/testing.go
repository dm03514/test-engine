package actions

import (
	"context"
	"github.com/dm03514/test-engine/results"
)

// DummyAction doesn't do anything!
type DummyAction struct {
	ReturnResult results.Result
}

// Execute a noop action
func (da DummyAction) Execute(context.Context, results.Results) (results.Result, error) {
	return da.ReturnResult, nil
}

// StubResult is a configurable result for testing
type StubResult struct{}

// Error returns a stub result error
func (sr StubResult) Error() error {
	return nil
}

// ValueOfProperty returns a value and an error for testing
func (sr StubResult) ValueOfProperty(property string) (results.Value, error) {
	return nil, nil
}
