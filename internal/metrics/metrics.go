package metrics

import (
	"strconv"

	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics interface {
	RecordHTTPServerRequest(method, route string, surface constants.Surface, status int, durationMs int64)
}

type metricsImpl struct {
	httpSrv *prometheus.HistogramVec
}

// Default buckets ( in ms)
var defaultBuckets = []float64{5, 10, 50, 100, 200, 500, 1000, 2000, 3000}

func NewMetrics() Metrics {
	httpSrv := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_server_request_ms",
		Help: "http server request duration",
		Buckets: defaultBuckets,
	}, []string{"method", "route", "surface", "status"})
	return &metricsImpl{
		httpSrv: httpSrv,
	}
}

func (m *metricsImpl) RecordHTTPServerRequest(method, route string, surface constants.Surface, status int, durationMs int64) {
	m.httpSrv.WithLabelValues(method, route, string(surface), strconv.Itoa(status)).Observe(float64(durationMs))
}