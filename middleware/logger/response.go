package logger

import (
	"net/http"

	"github.com/Compogo/compogo/logger"
)

// response wraps http.ResponseWriter to capture the response body for logging.
type response struct {
	http.ResponseWriter

	body []byte
}

// Write captures the written data while delegating to the underlying ResponseWriter.
func (response *response) Write(body []byte) (int, error) {
	response.body = append(response.body, body...)

	return response.ResponseWriter.Write(body)
}

// Response is middleware that logs HTTP response bodies at DEBUG level.
// Useful for debugging API responses in development environments.
type Response struct {
	logger logger.Logger
}

// NewResponse creates a new Response logging middleware.
func NewResponse(logger logger.Logger) *Response {
	return &Response{
		logger: logger.GetLogger("http.server.middleware.response"),
	}
}

// Middleware implements the http.Middleware interface.
// It captures the response body and logs it after the handler completes.
func (r *Response) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		responseWriter := &response{ResponseWriter: writer}

		next.ServeHTTP(responseWriter, request)

		r.logger.Debugf("path - '%s', body - '%s'", request.URL.String(), string(responseWriter.body))
	})
}
