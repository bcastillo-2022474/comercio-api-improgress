// Package metrics exposes Prometheus HTTP metrics for the server.
package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// unmatchedRoute labels requests that no mux pattern matched, keeping cardinality bounded.
const unmatchedRoute = "unmatched"

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests by method, matched route pattern and status code.",
	}, []string{"method", "route", "status"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds by method and matched route pattern.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route"})
)

// Handler serves the default registry, which already includes the Go runtime
// and process collectors.
func Handler() http.Handler {
	return promhttp.Handler()
}

// Middleware records request count and latency. The route label is the matched
// ServeMux pattern, read after the inner handler runs because the mux sets it
// on the request.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		route := r.Pattern
		if route == "" {
			route = unmatchedRoute
		}
		requestsTotal.WithLabelValues(r.Method, route, strconv.Itoa(rec.status)).Inc()
		requestDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Unwrap lets http.ResponseController reach the underlying writer.
func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}
