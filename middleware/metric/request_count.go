package metric

import (
	"net/http"
	"strconv"

	"github.com/Compogo/compogo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// responseStatusCode wraps http.ResponseWriter to capture the status code.
type responseStatusCode struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before delegating.
func (writer *responseStatusCode) WriteHeader(code int) {
	writer.statusCode = code
	writer.ResponseWriter.WriteHeader(code)
}

// RequestCount is middleware that counts HTTP requests by status code and endpoint.
// It exports a counter metric: compogo_http_server_requests_total{app, code, endpoint}
type RequestCount struct {
	counter *prometheus.CounterVec
}

// NewRequestCount creates a new RequestCount middleware.
// The appConfig provides the application name for the metric label.
func NewRequestCount(appConfig *compogo.Config) *RequestCount {
	return &RequestCount{
		counter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: compogo.MetricNamePrefix + "http_server_requests_total",
			Help: "Number of HTTP requests",
			ConstLabels: map[string]string{
				compogo.MetricAppNameFieldName: appConfig.Name,
			},
		}, []string{CodeFieldName, EndpointFieldName}),
	}
}

// Middleware implements the http.Middleware interface.
// It increments the counter after the request completes.
func (middleware *RequestCount) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingResponseWriter := &responseStatusCode{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(loggingResponseWriter, r)

		middleware.counter.With(prometheus.Labels{
			CodeFieldName:     strconv.Itoa(loggingResponseWriter.statusCode),
			EndpointFieldName: r.URL.Path,
		}).Inc()
	})
}
