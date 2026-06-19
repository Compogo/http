package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Compogo/compogo"
	"github.com/Compogo/runner"
)

// Server представляет HTTP-сервер с поддержкой graceful shutdown.
// Реализует интерфейс runner.Process для интеграции с Runner'ом.
type Server interface {
	runner.Process
	SetRouter(router Router)
}

type server struct {
	config *Config
	server *http.Server
	logger compogo.Logger
}

// NewServer создаёт новый HTTP-сервер.
// Принимает конфигурацию и логгер.
//
// Пример:
//
//	server := NewServer(config, logger)
//	server.SetRouter(router)
func NewServer(config *Config, logger compogo.Logger) Server {
	return &server{
		config: config,
		logger: logger.GetLogger("server.http"),
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", config.Interface, config.Port),
		},
	}
}

// Close останавливает HTTP-сервер с таймаутом.
// Реализует интерфейс io.Closer.
// Использует ShutdownTimeout из конфигурации.
func (server *server) Close() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), server.config.ShutdownTimeout)
	defer cancelFunc()

	server.logger.Info("shutdown")

	return server.server.Shutdown(ctx)
}

// Process запускает HTTP-сервер.
// Реализует интерфейс runner.Process.
func (server *server) Process(_ context.Context) error {
	return server.ListenAndServe()
}

// ListenAndServe запускает HTTP-сервер и обрабатывает ошибки.
// Игнорирует http.ErrServerClosed (штатное завершение).
func (server *server) ListenAndServe() error {
	server.logger.Infof("interface - %s, port - %d", server.config.Interface, server.config.Port)

	if err := server.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("[http.server] serve failed: %w", err)
	}

	return nil
}

func (server *server) Name() string {
	return "server.http"
}

func (server *server) SetRouter(router Router) {
	server.server.Handler = router
}
