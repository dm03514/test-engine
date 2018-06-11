package observables

import "context"

// ObservableEvent is emitted by the Observable
type ObservableEvent interface{}

// Observable can be watched for events by the engine
type Observable interface {
	RunUntilContextDone(context.Context) <-chan ObservableEvent
	Name() string
}
