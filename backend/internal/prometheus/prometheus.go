package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Prometheus struct {
	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
	HttpErrorsTotal     *prometheus.CounterVec

	DBQueryDuration *prometheus.HistogramVec
	DBErrorsTotal   *prometheus.CounterVec
}

func NewPrometheus() *Prometheus {
	p := &Prometheus{
		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),

		HttpErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_errors_total",
				Help: "Total number of HTTP errors",
			},
			[]string{"method", "path", "status"},
		),

		DBQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Duration of DB queries in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"query_type"},
		),

		DBErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_errors_total",
				Help: "Total number of DB errors",
			},
			[]string{"query_type"},
		),
	}

	// Регистрация метрик
	prometheus.MustRegister(
		p.HttpRequestsTotal,
		p.HttpRequestDuration,
		p.HttpErrorsTotal,
		p.DBQueryDuration,
		p.DBErrorsTotal,
	)

	return p
}
