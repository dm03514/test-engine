package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	pass string = "pass"
	fail string = "fail"
)

func errToResult(err error) string {
	if err == nil {
		return pass
	}

	return fail
}

type HistogramVecDurationRecorder struct {
	*prometheus.HistogramVec
}

func (h *HistogramVecDurationRecorder) Record(d time.Duration, err error) {
	h.HistogramVec.With(prometheus.Labels{"result": errToResult(err)}).Observe(d.Seconds())
}

/*
func New(t engine.Test, opts ...engine.Option) (*engine.Engine, error) {

	stateDuration := HistogramVecDurationRecorder{
		HistogramVec: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "state_duration",
				Help: "Duration of an individual state result:pass|fail",
			},
			[]string{"result"},
		),
	}
	if err := prometheus.Register(stateDuration); err != nil {
		return nil, err
	}

	testDuration := HistogramVecDurationRecorder{
		HistogramVec: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "test_duration",
				Help: "Duration of a complete test with result:pass|fail",
			},
			[]string{"result"},
		),
	}
	if err := prometheus.Register(testDuration); err != nil {
		return nil, err
	}

	// prepend opts to opts
	opts = append(
		[]engine.Option{
			engine.OptionRecordStateDuration(stateDuration.Record),
			engine.OptionRecordTestDuration(testDuration.Record),
		},
		opts...,
	)

	return engine.New(
		t,
		opts...,
	)
}
*/
