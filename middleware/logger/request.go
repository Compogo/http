package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/http/helper"
)

// Request is middleware that logs HTTP request bodies at DEBUG level.
// It reads the entire body, logs it, and restores it for the next handler.
type Request struct {
	logger logger.Logger
}

// NewRequest creates a new Request logging middleware.
func NewRequest(logger logger.Logger) *Request {
	return &Request{
		logger: logger.GetLogger("http.server.middleware.request"),
	}
}

// Middleware implements the http.Middleware interface.
// It reads the request body, logs it, and restores it for subsequent handlers.
// If body reading fails, it returns 400 Bad Request.
func (r *Request) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil && err != io.EOF {
			err = fmt.Errorf("read body failed: %w", err)
			r.logger.Error(err)
			helper.WriteError(writer, request, err.Error(), http.StatusBadRequest)
			return
		}

		request.Body = io.NopCloser(bytes.NewReader(body))

		r.logger.Debugf("path - '%s', body - '%s'", request.URL.String(), string(body))

		next.ServeHTTP(writer, request)
	})
}
