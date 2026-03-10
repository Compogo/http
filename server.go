package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/runner"
)

// Server defines the interface for an HTTP server component.
// It integrates with runner.Runner for lifecycle management
// and supports graceful shutdown via io.Closer.
type Server interface {
	// Closer io.Closer provides graceful shutdown functionality.
	// It waits for active requests to complete up to ShutdownTimeout.
	io.Closer

	// Process runner.Process allows the server to be run as a task in the runner.
	// The server will block until the context is canceled or an error occurs.
	runner.Process

	// SetRouter attaches a Router to the server.
	// Must be called before starting the server.
	SetRouter(router Router)
}

type server struct {
	config *Config
	server *http.Server
	logger logger.Logger
}

// NewServer creates a new HTTP server instance with the given configuration.
// The server is not started until it's passed to runner.Runner.
// The logger is automatically namespaced with "server.http" for better log filtering.
func NewServer(config *Config, logger logger.Logger) Server {
	return &server{
		config: config,
		logger: logger.GetLogger("server.http"),
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", config.Interface, config.Port),
		},
	}
}

func (server *server) Close() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), server.config.ShutdownTimeout)
	defer cancelFunc()

	server.logger.Info("shutdown")

	return server.server.Shutdown(ctx)
}

func (server *server) Process(ctx context.Context) (err error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	server.logger.Infof("interface - %s, port - %d", server.config.Interface, server.config.Port)

	if err = server.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("[http.server] serve failed: %w", err)
	}

	return nil
}

func (server *server) SetRouter(router Router) {
	server.server.Handler = router
}
