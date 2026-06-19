package basic

import (
	"net/http"

	"github.com/Compogo/http_server/helper"
)

// Auth реализует HTTP Basic Authentication middleware.
// Проверяет учетные данные из заголовка Authorization.
type Auth struct {
	config *Config
}

// NewAuth создаёт новый middleware для Basic Auth.
func NewAuth(config *Config) *Auth {
	return &Auth{config: config}
}

// Middleware реализует интерфейс http_server.Middleware.
// Проверяет логин и пароль из заголовка Authorization.
// В случае неудачи возвращает 401 Unauthorized.
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
