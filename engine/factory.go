package engine

type Factory interface {
	New(t Test, opts ...Option) (*Engine, error)
}

type DefaultFactory struct{}

func (d DefaultFactory) New(t Test, opts ...Option) (*Engine, error) {
	e := &Engine{
		Test:                t,
		recordStateDuration: NoopDurationRecorder,
		recordTestDuration:  NoopDurationRecorder,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e, nil
}

func NewDefaultFactory() DefaultFactory {
	return DefaultFactory{}
}
