package main

import (
	"flag"
	"fmt"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/engine"
	ep "github.com/dm03514/test-engine/engine/prometheus"
	"github.com/dm03514/test-engine/transcons"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"os"
)

type HTTPExecutor interface {
	ListenAndServe()
	RegisterHandlers()
}

type NewServer func(loaders engine.Loaders, opts ...engine.HTTPOpt) (HTTPExecutor, error)

func defaultServer(testsDir string) (HTTPExecutor, error) {
	ar, err := actions.NewActionRegistry()
	if err != nil {
		return nil, err
	}

	tcr, err := transcons.NewTransConsRegistry()
	if err != nil {
		return nil, err
	}

	loader, err := engine.NewFileLoader(testsDir, ar, tcr, engine.NewDefaultFactory())
	if err != nil {
		return nil, err
	}

	return engine.NewHTTPExecutor(
		engine.NewLoaders(
			loader,
		),
	)
}

func prometheusServer(testsDir string) (HTTPExecutor, error) {
	stateDuration := ep.HistogramVecStateDurationRecorder{
		HistogramVec: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "state_duration",
				Help: "Duration of an individual state result:pass|fail",
			},
			[]string{"result", "state_name", "test_name"},
		),
	}

	testDuration := ep.HistogramVecTestDurationRecorder{
		prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "test_duration_seconds",
				Help: "Duration of a complete test with result:pass|fail",
			},
			[]string{"result", "test_name"},
		),
	}

	prometheus.MustRegister(
		testDuration,
		stateDuration,
	)

	ar, err := actions.NewActionRegistry()
	if err != nil {
		return nil, err
	}

	tcr, err := transcons.NewTransConsRegistry()
	if err != nil {
		return nil, err
	}

	loader, err := engine.NewFileLoader(
		testsDir,
		ar,
		tcr,
		engine.NewDefaultFactory(
			engine.OptionRecordStateDuration(stateDuration.Record),
			engine.OptionRecordTestDuration(testDuration.Record),
		),
	)
	if err != nil {
		return nil, err
	}

	return ep.NewPrometheusHTTPExecutor(
		engine.NewLoaders(
			loader,
		),
	)

	return nil, nil
}

func main() {
	var testsDir = flag.String("testDir", "", "Path to serve tests from")
	var metrics = flag.String("metrics", "default", "Metrics to use: default|prometheus")
	flag.Parse()

	log.Infof("testsDir: %q", *testsDir)
	log.Infof("metrics: %q", *metrics)
	var s HTTPExecutor
	var err error

	switch *metrics {
	case "default":
		s, err = defaultServer(*testsDir)
	case "prometheus":
		s, err = prometheusServer(*testsDir)
	default:
		panic(fmt.Errorf("%q metric not supported", *metrics))
	}
	fmt.Print(s)

	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	// http.Handle("/metrics", promhttp.Handler())
	// log.Fatal(http.ListenAndServe(":8080", nil))
	/*
		http.HandleFunc("/execute", s.Execute)
		http.Handle("/metrics", promhttp.Handler())
	*/
	s.RegisterHandlers()
	s.ListenAndServe()
}
