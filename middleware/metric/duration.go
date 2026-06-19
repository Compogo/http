package metric

import (
	"net/http"

	"github.com/Compogo/compogo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Duration — middleware для сбора метрики длительности HTTP-запросов.
// Измеряет время выполнения каждого запроса в секундах.
//
// Метрика: compogo_http_server_duration_seconds{app="myapp", endpoint="/api/v1/users"}
type Duration struct {
	counter *prometheus.HistogramVec
}

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

func (middleware *Duration) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(middleware.counter.With(prometheus.Labels{
			EndpointFieldName: r.URL.Path,
		}))
		defer timer.ObserveDuration()

		next.ServeHTTP(w, r)
	})
}
