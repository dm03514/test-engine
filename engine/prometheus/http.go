package prometheus

import (
	"github.com/dm03514/test-engine/engine"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// HTTPExecutor wraps the vanilla HTTP executor
type HTTPExecutor struct {
	engine.HttpExecutor
}

// RegisterHandlers adds prometheus metrics endpoint
func (p HTTPExecutor) RegisterHandlers() {
	http.Handle("/metrics", promhttp.Handler())

	p.HttpExecutor.RegisterHandlers()
}

// NewHTTPExecutor initializes prometheus HTTP executor with options
func NewHTTPExecutor(loaders engine.Loaders, opts ...engine.HTTPOpt) (HTTPExecutor, error) {
	e, err := engine.NewHTTPExecutor(
		loaders,
		opts...,
	)
	if err != nil {
		return HTTPExecutor{}, nil
	}
	return HTTPExecutor{
		HttpExecutor: e,
	}, nil
}
