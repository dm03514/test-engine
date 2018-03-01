package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
	log "github.com/sirupsen/logrus"
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
	log.Infof("prometheus.Record(%s, %+v, errToResult(%q), %+v", d, err, errToResult(err), h.HistogramVec)
}
