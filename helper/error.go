package helper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/thoas/go-funk"
)

// MimeTypeNOSNIFF — значение заголовка X-Content-Type-Options
// для предотвращения MIME-сниффинга.
const MimeTypeNOSNIFF = "nosniff"

// Error — структура ошибки для JSON-ответов.
type Error struct {
	error string
}

// NewError создаёт новую ошибку.
func NewError(error string) *Error {
	return &Error{error: error}
}

// Error возвращает сообщение об ошибке.
// Реализует интерфейс error.
func (e *Error) Error() string {
	return e.error
}

// MarshalJSON сериализует ошибку в JSON.
// Формат: {"error": "сообщение"}
func (e *Error) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"error\":\"%s\"}", e.Error())), nil
}

// JSONError отправляет ошибку в формате JSON.
// Устанавливает заголовки Content-Type и X-Content-Type-Options.
//
// Пример:
//
//	helper.JSONError(w, helper.NewError("user not found"), http.StatusNotFound)
func JSONError(w http.ResponseWriter, err *Error, code int) {
	w.Header().Set(headers.ContentType, mimetype.ApplicationJSON)
	w.Header().Set(headers.XContentTypeOptions, MimeTypeNOSNIFF)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(err)
}

// WriteError отправляет ошибку в формате JSON или plain text в зависимости
// от заголовков Accept и Content-Type запроса.
//
// Если клиент ожидает JSON (application/json), отправляет JSON-ответ.
// Иначе отправляет plain text через http.Error.
//
// Пример:
//
//	helper.WriteError(w, r, "user not found", http.StatusNotFound)
func WriteError(w http.ResponseWriter, r *http.Request, err string, code int) {
	acceptTypes, existAcceptTypes := r.Header[headers.Accept]
	contentTypes, existContentTypes := r.Header[headers.ContentType]

	if (existAcceptTypes && funk.ContainsString(acceptTypes, mimetype.ApplicationJSON)) || (existContentTypes && funk.ContainsString(contentTypes, mimetype.ApplicationJSON)) {
		JSONError(w, NewError(err), code)
		return
	}

	http.Error(w, err, code)
}
