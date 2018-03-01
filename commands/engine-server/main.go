package main

import (
	"flag"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/engine"
	ep "github.com/dm03514/test-engine/engine/prometheus"
	"github.com/dm03514/test-engine/transcons"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"os"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"fmt"
)

type HTTPExecutor interface {
	Execute(http.ResponseWriter, *http.Request)
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
	/*
	stateDuration := ep.HistogramVecDurationRecorder{
		HistogramVec: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "state_duration",
				Help: "Duration of an individual state result:pass|fail",
			},
			[]string{"result"},
		),
	}
	*/

	testDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "test_duration_seconds",
			Help: "Duration of a complete test with result:pass|fail",
		},
		[]string{"result"},
	)

	fmt.Printf("REGISTERING\n")
	prometheus.MustRegister(testDuration)

	/*
	t, err := time.ParseDuration("1s")
	if err != nil {
		panic(err)
	}
	*/
	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("OBSERVING\n")
	testDuration.With(prometheus.Labels{"result": "pass"}).Observe(1)
	testDuration.With(prometheus.Labels{"result": "pass"}).Observe(1)
	testDuration.With(prometheus.Labels{"result": "pass"}).Observe(1)
	// testDuration.Record(t, nil)

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
			// engine.OptionRecordStateDuration(stateDuration.Record),
			// engine.OptionRecordTestDuration(testDuration.Record),
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
	log.Fatal(http.ListenAndServe(":8080", nil))
	/*
	// s.RegisterHandlers()
	http.HandleFunc("/execute", s.Execute)
	http.Handle("/metrics", promhttp.Handler())
	s.ListenAndServe()
	*/
}
