package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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

// HistogramVecStateDurationRecorder can record state durations
type HistogramVecStateDurationRecorder struct {
	*prometheus.HistogramVec
}

// Record a metric
func (h *HistogramVecStateDurationRecorder) Record(sn string, tn string, d time.Duration, err error) {
	log.Infof("prometheus.Record(%s, %+v, errToResult(%q), %q, %q, %+v",
		d, err, errToResult(err), sn, tn, h.HistogramVec)
	h.HistogramVec.With(prometheus.Labels{
		"result":     errToResult(err),
		"state_name": sn,
		"test_name":  tn,
	}).Observe(d.Seconds())
}

// HistogramVecTestDurationRecorder can record test durations to prometheus!
type HistogramVecTestDurationRecorder struct {
	*prometheus.HistogramVec
}

// Record a histogram vector to prometheus
func (h *HistogramVecTestDurationRecorder) Record(tn string, d time.Duration, err error) {
	log.Infof("prometheus.Record(%s, %+v, errToResult(%q), %q, %+v",
		d, err, errToResult(err), tn, h.HistogramVec)
	h.HistogramVec.With(prometheus.Labels{
		"result":    errToResult(err),
		"test_name": tn,
	}).Observe(d.Seconds())
}
