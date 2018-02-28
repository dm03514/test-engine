package engine

type Factory interface {
	New(t Test, opts ...Option) (*Engine, error)
}

type DefaultFactory struct {
	extraOpts []Option
}

func (d DefaultFactory) New(t Test, opts ...Option) (*Engine, error) {
	e := &Engine{
		Test:                t,
		recordStateDuration: NoopDurationRecorder,
		recordTestDuration:  NoopDurationRecorder,
	}

	for _, opt := range d.extraOpts {
		opt(e)
	}

	for _, opt := range opts {
		opt(e)
	}

	return e, nil
}

func NewDefaultFactory(eo ...Option) DefaultFactory {
	return DefaultFactory{
		extraOpts: eo,
	}
}
