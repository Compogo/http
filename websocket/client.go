package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Compogo/compogo"
	"github.com/Compogo/types/emitter"
	"github.com/gorilla/websocket"
)

// Client представляет WebSocket-клиента (одно соединение).
// Обеспечивает:
//   - Двусторонний обмен сообщениями
//   - Автоматические ping/pong для поддержания соединения
//   - Буферизацию исходящих сообщений
//   - Graceful shutdown через контекст
type Client struct {
	conn   *websocket.Conn
	config *Config

	events chan *Event
	ticker *time.Ticker

	onMessage emitter.Emitter[*Event]
	logger    compogo.Logger
}

// NewClient создаёт нового WebSocket-клиента.
func NewClient(conn *websocket.Conn, config *Config, onMessage emitter.Emitter[*Event], logger compogo.Logger) *Client {
	return &Client{
		conn:      conn,
		config:    config,
		events:    make(chan *Event, config.ClientEventBufferSize),
		ticker:    time.NewTicker(config.PingInterval),
		onMessage: onMessage,
		logger:    logger,
	}
}

// Send отправляет событие клиенту.
// Если буфер исходящих сообщений заполнен, возвращает MessageChanFullError.
func (c *Client) Send(event *Event) error {
	if event.Timestamp == nil {
		event.Timestamp = NewTimestamp()
	}

	select {
	case c.events <- event:
		return nil
	default:
		return MessageChanFullError
	}
}

// Process запускает основной цикл обработки клиента.
// Блокирует выполнение до завершения работы клиента.
// Обрабатывает входящие и исходящие сообщения, ping/pong.
func (c *Client) Process(ctx context.Context) error {
	mainCtx, mainCancel := context.WithCancel(ctx)
	defer c.conn.Close()
	defer close(c.events)
	defer mainCancel()
	defer c.ticker.Stop()

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(c.config.PingTimeout))
	})

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		readCtx, readCancel := context.WithCancel(mainCtx)
		defer mainCancel()
		defer readCancel()

		if err := c.conn.SetReadDeadline(time.Now().Add(c.config.PingTimeout)); err != nil {
			c.logger.Errorf("websocket: failed to set read deadline: %s", err.Error())
			return
		}

		for {
			t, message, err := c.conn.ReadMessage()
			if err != nil && websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				return
			}

			if err != nil {
				c.logger.Errorf("websocket: failed to read message: %s", err.Error())
				return
			}

			if t == websocket.TextMessage {
				event := &Event{}
				if err := json.Unmarshal(message, event); err != nil {
					c.logger.Errorf("websocket: failed to unmarshal event: %s", err.Error())
					continue
				}

				c.onMessage.Emit(readCtx, event)
			}

			select {
			case <-readCtx.Done():
				return
			default:
				break
			}
		}
	})

	wg.Go(func() {
		writeCtx, writeCancel := context.WithCancel(mainCtx)
		defer mainCancel()
		defer writeCancel()

		for {
			select {
			case <-c.ticker.C:
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout)); err != nil {
					c.logger.Errorf("failed to set write deadline: %s", err.Error())
					return
				}

				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					c.logger.Errorf("failed to write ping message: %s", err.Error())
					return
				}
			case msg, ok := <-c.events:
				if !ok {
					c.logger.Error(MessageChanClosedError)
					return
				}

				if err := c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout)); err != nil {
					c.logger.Errorf("failed to set write deadline: %s", err.Error())
					return
				}

				jsonBody, err := json.Marshal(msg)
				if err != nil {
					c.logger.Error("json marshal:", err)
					continue
				}

				if err = c.conn.WriteMessage(websocket.TextMessage, jsonBody); err != nil {
					c.logger.Errorf("failed to write message: %s", err.Error())
					return
				}
			case <-writeCtx.Done():
				return
			}
		}
	})

	wg.Wait()
	return nil
}
