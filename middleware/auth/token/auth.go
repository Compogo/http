package token

import (
	"net/http"

	"github.com/Compogo/http/helper"
)

// Auth implements HTTP Token Authentication middleware.
// It validates incoming requests against configured tokens.
type Auth struct {
	config *Config
}

// NewAuth creates a new Auth middleware instance with the given configuration.
func NewAuth(config *Config) *Auth {
	return &Auth{config: config}
}

// Middleware implements the http.Middleware interface.
// It checks for a valid token in the configured header:
//   - If token missing or invalid, returns 401 Unauthorized
//   - If valid, passes the request to the next handler
func (auth *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !auth.config.Tokens.Contains(request.Header.Get(auth.config.HeaderName)) {
			helper.WriteError(writer, request, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
