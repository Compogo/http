package metric

import (
	"net/http"

	"github.com/Compogo/compogo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Duration is middleware that measures HTTP request duration.
// It exports a histogram metric: compogo_http_server_duration_seconds{app, endpoint}
type Duration struct {
	counter *prometheus.HistogramVec
}

// NewDuration creates a new Duration middleware.
// The appConfig provides the application name for the metric label.
func NewDuration(appConfig *compogo.Config) *Duration {
	return &Duration{
		counter: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: compogo.MetricNamePrefix + "http_server_duration_seconds",
			Help: "Duration of HTTP requests",
			ConstLabels: map[string]string{
				compogo.MetricAppNameFieldName: appConfig.Name,
			},
		}, []string{EndpointFieldName}),
	}
}

// Middleware implements the http.Middleware interface.
// It measures the duration of the request and records it in the histogram.
func (middleware *Duration) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(middleware.counter.With(prometheus.Labels{
			EndpointFieldName: r.URL.Path,
		}))
		defer timer.ObserveDuration()

		next.ServeHTTP(w, r)
	})
}
