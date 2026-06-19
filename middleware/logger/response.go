package logger

import (
	"net/http"

	"github.com/Compogo/compogo"
)

// response — обёртка над http.ResponseWriter для захвата тела ответа.
type response struct {
	http.ResponseWriter

	body []byte
}

func (response *response) Write(body []byte) (int, error) {
	response.body = append(response.body, body...)

	return response.ResponseWriter.Write(body)
}

// Response — middleware для логирования исходящих HTTP-ответов.
// Логирует тело ответа.
type Response struct {
	logger compogo.Logger
}

func NewResponse(logger compogo.Logger) *Response {
	return &Response{
		logger: logger.GetLogger("http").GetLogger("server").GetLogger("middleware").GetLogger("response"),
	}
}

func (r *Response) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		responseWriter := &response{ResponseWriter: writer}

		next.ServeHTTP(responseWriter, request)

		r.logger.Debugf("path - '%s', body - '%s'", request.URL.String(), string(responseWriter.body))
	})
}
