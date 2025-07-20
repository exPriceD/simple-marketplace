package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type PromMetrics struct {
	reqCount   *prometheus.CounterVec
	reqLatency *prometheus.HistogramVec

	ucCount   *prometheus.CounterVec
	ucLatency *prometheus.HistogramVec
	registry  *prometheus.Registry
}

func New() *PromMetrics {
	r := prometheus.NewRegistry()

	pm := &PromMetrics{
		reqCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		}, []string{"method", "path", "status"}),

		reqLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path"}),

		ucCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "usecase_calls_total",
			Help: "Use case calls",
		}, []string{"op", "status"}),

		ucLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "usecase_duration_seconds",
			Help:    "Use case execution latency",
			Buckets: prometheus.DefBuckets,
		}, []string{"op"}),

		registry: r,
	}

	r.MustRegister(pm.reqCount, pm.reqLatency, pm.ucCount, pm.ucLatency)
	return pm
}

func (m *PromMetrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

func (m *PromMetrics) IncHTTPRequests(method, path, status string) {
	m.reqCount.WithLabelValues(method, path, status).Inc()
}

func (m *PromMetrics) ObserveHTTPDuration(method, path string, seconds float64) {
	m.reqLatency.WithLabelValues(method, path).Observe(seconds)
}

func (m *PromMetrics) IncUseCaseCalls(op, status string) {
	m.ucCount.WithLabelValues(op, status).Inc()
}

func (m *PromMetrics) ObserveUseCaseDuration(op string, seconds float64) {
	m.ucLatency.WithLabelValues(op).Observe(seconds)
}
