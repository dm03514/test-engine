package prometheus

import "github.com/dm03514/test-engine/engine"

func New(t engine.Test, opts ...engine.Option) (*engine.Engine, error) {
	return engine.New(t, opts...)
}
