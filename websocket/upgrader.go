package websocket

import (
	"net/http"
	"net/url"

	"github.com/go-http-utils/headers"
	"github.com/gorilla/websocket"
)

// Upgrader — обёртка над websocket.Upgrader с поддержкой проверки Origin
type Upgrader struct {
	*websocket.Upgrader
	config *Config
}

// NewUpgrader создаёт новый Upgrader с конфигурацией из Config.
func NewUpgrader(config *Config) *Upgrader {
	up := &Upgrader{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  config.ReadBufferSize,
			WriteBufferSize: config.WriteBufferSize,
		},
		config: config,
	}

	up.Upgrader.CheckOrigin = up.CheckOrigin

	return up
}

// CheckOrigin проверяет, разрешён ли Origin запроса.
// Поддерживает:
//   - Wildcard "*" — разрешены все origins
//   - Совпадение с хостом запроса (same-origin)
//   - Список разрешённых origins из конфигурации
//
// Реализует интерфейс websocket.Upgrader.CheckOrigin.
func (up *Upgrader) CheckOrigin(r *http.Request) bool {
	if up.config.Origins.Contains(AnyOrigin) {
		return true
	}

	origin := r.Header.Get(headers.Origin)
	if origin == "" {
		return false
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if u.Host == r.Host {
		return true
	}

	return up.config.Origins.Contains(u.Host)
}
