package prometheus

import (
	"github.com/dm03514/test-engine/engine"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type PrometheusHTTPExecutor struct {
	engine.HttpExecutor
}

func (p PrometheusHTTPExecutor) RegisterHandlers() {
	http.Handle("/metrics", promhttp.Handler())

	p.HttpExecutor.RegisterHandlers()
}

func NewPrometheusHTTPExecutor(loaders engine.Loaders, opts ...engine.HTTPOpt) (PrometheusHTTPExecutor, error) {
	e, err := engine.NewHTTPExecutor(
		loaders,
		opts...,
	)
	if err != nil {
		return PrometheusHTTPExecutor{}, nil
	}
	return PrometheusHTTPExecutor{
		HttpExecutor: e,
	}, nil
}
