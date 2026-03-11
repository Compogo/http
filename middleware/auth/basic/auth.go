package basic

import (
	"net/http"

	"github.com/Compogo/http/helper"
)

// Auth implements HTTP Basic Authentication middleware.
// It validates incoming requests against configured credentials.
type Auth struct {
	config *Config
}

// NewAuth creates a new Auth middleware instance with the given configuration.
func NewAuth(config *Config) *Auth {
	return &Auth{config: config}
}

// Middleware implements the http.Middleware interface.
// It checks for valid Basic Authentication credentials:
//   - If no credentials or invalid credentials, returns 401 Unauthorized
//   - If valid, passes the request to the next handler
func (auth *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userName, password, ok := request.BasicAuth()
		if !ok {
			helper.WriteError(writer, request, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		cred, err := auth.config.Creds.Get(userName)
		if err != nil || cred.Password != password {
			helper.WriteError(writer, request, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
