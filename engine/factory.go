package engine

import "github.com/dm03514/test-engine/observables"

// Factory creates new engines!
type Factory interface {
	New(t Test, opts ...Option) (*Engine, error)
}

// DefaultFactory contains options and ability to create engines
type DefaultFactory struct {
	extraOpts []Option
}

// New initializes an engine using this factory, with the given options
// and decorates the test
func (d DefaultFactory) New(t Test, opts ...Option) (*Engine, error) {
	e := &Engine{
		Test:                t,
		eventChans:          make(map[string]<-chan observables.ObservableEvent),
		recordStateDuration: NoopStateDurationRecorder,
		recordTestDuration:  NoopTestDurationRecorder,
	}

	for _, opt := range d.extraOpts {
		opt(e)
	}

	for _, opt := range opts {
		opt(e)
	}

	return e, nil
}

// NewDefaultFactory initializes a factory
func NewDefaultFactory(eo ...Option) DefaultFactory {
	return DefaultFactory{
		extraOpts: eo,
	}
}
