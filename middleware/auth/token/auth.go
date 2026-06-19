package token

import (
	"net/http"

	"github.com/Compogo/http_server/helper"
)

// Auth реализует Token-based Authentication middleware.
// Проверяет наличие токена в заголовке запроса.
type Auth struct {
	config *Config
}

// NewAuth создаёт новый middleware для Token Auth.
func NewAuth(config *Config) *Auth {
	return &Auth{config: config}
}

// Middleware реализует интерфейс http_server.Middleware.
// Проверяет наличие токена в заголовке.
// В случае неудачи возвращает 401 Unauthorized.
func (auth *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !auth.config.Tokens.Contains(request.Header.Get(auth.config.HeaderName)) {
			helper.WriteError(writer, request, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
