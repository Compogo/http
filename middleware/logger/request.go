package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Compogo/compogo"
	"github.com/Compogo/http_server/helper"
)

// Request — middleware для логирования входящих HTTP-запросов.
// Логирует путь и тело запроса.
type Request struct {
	logger compogo.Logger
}

func NewRequest(logger compogo.Logger) *Request {
	return &Request{
		logger: logger.GetLogger("http").GetLogger("server").GetLogger("middleware").GetLogger("request"),
	}
}

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
