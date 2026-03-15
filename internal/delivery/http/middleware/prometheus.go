package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de requisições HTTP.",
		},
		[]string{"path", "method", "code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP em segundos.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

// responseWriter é um wrapper para http.ResponseWriter para capturar o status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// PrometheusMiddleware é o middleware que coleta as métricas.
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := NewResponseWriter(w)

		// Chama o próximo handler na cadeia
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Tenta obter a rota do mux para evitar cardinalidade alta com IDs
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		if path == "" {
			path = r.URL.Path
		}

		// Coleta as métricas
		httpRequestDuration.WithLabelValues(path, r.Method).Observe(duration.Seconds())
		httpRequestsTotal.WithLabelValues(path, r.Method, strconv.Itoa(rw.statusCode)).Inc()
	})
}
