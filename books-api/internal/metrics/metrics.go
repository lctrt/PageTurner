package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pageturner_http_requests_total",
			Help: "Total HTTP requests by method, route and status.",
		},
		[]string{"method", "route", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pageturner_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	GoodreadsImports = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pageturner_goodreads_imports_total",
			Help: "Total Goodreads import attempts by result.",
		},
		[]string{"result"},
	)

	PagesReadTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pageturner_pages_read_total",
			Help: "Total pages read across all users since process start.",
		},
	)
)

func Register(reg prometheus.Registerer) {
	reg.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
		GoodreadsImports,
		PagesReadTotal,
	)
}
