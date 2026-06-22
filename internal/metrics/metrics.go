package metrics

// NOTE: Add to go.mod: github.com/prometheus/client_golang v1.20.0

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"},
		[]string{"method", "path", "status"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration"},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(RequestsTotal, RequestDuration)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
