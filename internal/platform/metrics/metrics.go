package metrics

import "context"

type Metrics interface {
	IncHTTPRequests(method, path, status string)
	ObserveHTTPDuration(method, path string, seconds float64)

	IncUseCaseCalls(op, status string)
	ObserveUseCaseDuration(op string, seconds float64)
}

type NoopMetrics struct{}

func (NoopMetrics) IncHTTPRequests(_, _, _ string)             {}
func (NoopMetrics) ObserveHTTPDuration(_, _ string, _ float64) {}
func (NoopMetrics) IncUseCaseCalls(_, _ string)                {}
func (NoopMetrics) ObserveUseCaseDuration(_ string, _ float64) {}

type CtxMetrics interface {
	Metrics
}

func WithContext(ctx context.Context, m Metrics) context.Context { return ctx }
