package helper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/thoas/go-funk"
)

const (
	// MimeTypeNOSNIFF prevents browsers from MIME type sniffing
	MimeTypeNOSNIFF = "nosniff"
)

// Error represents a structured error response for JSON APIs.
// It implements both error and fmt.Stringer interfaces.
type Error struct {
	error string
}

// NewError creates a new Error with the given message.
func NewError(error string) *Error {
	return &Error{error: error}
}

// Error returns the error message, implementing the error interface.
func (e *Error) Error() string {
	return e.error
}

// MarshalJSON implements json.Marshaler for consistent JSON error format.
// Produces: {"error":"message"}
func (e *Error) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"error\":\"%s\"}", e.Error())), nil
}

// JSONError writes a JSON-formatted error response.
// It sets appropriate headers (Content-Type, X-Content-Type-Options)
// and writes the error with the given status code.
func JSONError(w http.ResponseWriter, err *Error, code int) {
	w.Header().Set(headers.ContentType, mimetype.ApplicationJSON)
	w.Header().Set(headers.XContentTypeOptions, MimeTypeNOSNIFF)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(err)
}

// WriteError intelligently writes an error response based on the request's
// Accept and Content-Type headers. If the client expects JSON, it responds
// with JSONError; otherwise, it falls back to plain text via http.Error.
func WriteError(w http.ResponseWriter, r *http.Request, err string, code int) {
	acceptTypes, existAcceptTypes := r.Header[headers.Accept]
	contentTypes, existContentTypes := r.Header[headers.ContentType]

	if (existAcceptTypes && funk.ContainsString(acceptTypes, mimetype.ApplicationJSON)) || (existContentTypes && funk.ContainsString(contentTypes, mimetype.ApplicationJSON)) {
		JSONError(w, NewError(err), code)
		return
	}

	http.Error(w, err, code)
}
