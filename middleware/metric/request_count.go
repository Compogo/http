package metric

import (
	"net/http"
	"strconv"

	"github.com/Compogo/compogo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// responseStatusCode — обёртка над http.ResponseWriter для захвата HTTP-кода ответа.
type responseStatusCode struct {
	http.ResponseWriter
	statusCode int
}

func (writer *responseStatusCode) WriteHeader(code int) {
	writer.statusCode = code
	writer.ResponseWriter.WriteHeader(code)
}

// RequestCount — middleware для сбора метрики количества HTTP-запросов.
// Считает количество запросов с разбивкой по эндпоинтам и кодам ответа.
//
// Метрика: compogo_http_server_requests_total{app="myapp", code="200", endpoint="/api/v1/users"}
type RequestCount struct {
	counter *prometheus.CounterVec
}

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
