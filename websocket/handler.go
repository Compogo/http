package websocket

import (
	"context"
	"net/http"
	"sync"

	"github.com/Compogo/compogo/closer"
	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/types/emitter"
	"github.com/Compogo/types/set"
)

type Handler struct {
	config   *Config
	upgrader *Upgrader

	rwm     sync.RWMutex
	clients set.Set[*Client]

	OnClientConnection    emitter.Emitter[*Client]
	OnClientDisconnection emitter.Emitter[*Client]
	OnMessage             emitter.Emitter[*Event]

	logger logger.Logger
	closer closer.Closer
}

func NewHandler(
	config *Config,
	upgrader *Upgrader,
	logger logger.Logger,
	closer closer.Closer,
) *Handler {
	return &Handler{
		config:                config,
		upgrader:              upgrader,
		OnClientConnection:    emitter.NewEmitter[*Client](),
		OnClientDisconnection: emitter.NewEmitter[*Client](),
		OnMessage:             emitter.NewEmitter[*Event](),
		logger:                logger.GetLogger("websocket"),
		closer:                closer,
	}
}

func (h *Handler) Send(event *Event) (err error) {
	h.rwm.RLock()
	defer h.rwm.RUnlock()

	event.Timestamp = NewTimestamp()

	for client := range h.clients {
		if err = client.Send(event); err != nil {
			h.logger.Error(err)
		}
	}

	return nil
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	conn, err := h.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		h.logger.Errorf("http -> websocket upgrade failed: %s", err.Error())
		return
	}

	client := NewClient(conn, h.config, h.OnMessage, h.logger)

	ctx := request.Context()

	h.addClient(ctx, client)
	defer h.removeClient(ctx, client)

	if err := client.Process(ctx); err != nil {
		h.logger.Errorf("process failed: %s", err.Error())
	}
}

func (h *Handler) addClient(ctx context.Context, client *Client) {
	h.rwm.Lock()
	defer h.rwm.Unlock()
	defer h.OnClientConnection.Emit(ctx, client)

	h.clients.Add(client)
}

func (h *Handler) removeClient(ctx context.Context, client *Client) {
	h.rwm.Lock()
	defer h.rwm.Unlock()
	defer h.OnClientDisconnection.Emit(ctx, client)

	h.clients.Remove(client)
}
