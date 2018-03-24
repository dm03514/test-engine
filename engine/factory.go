package engine

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
