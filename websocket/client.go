package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/types/emitter"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	config *Config

	events chan *Event
	ticker *time.Ticker

	onMessage emitter.Emitter[*Event]
	logger    logger.Logger
}

func NewClient(conn *websocket.Conn, config *Config, onMessage emitter.Emitter[*Event], logger logger.Logger) *Client {
	return &Client{
		conn:      conn,
		config:    config,
		events:    make(chan *Event, config.ClientEventBufferSize),
		ticker:    time.NewTicker(config.PingInterval),
		onMessage: onMessage,
		logger:    logger,
	}
}

func (c *Client) Send(event *Event) error {
	select {
	case c.events <- event:
		return nil
	default:
		return MessageChanFullError
	}
}

func (c *Client) Process(ctx context.Context) error {
	mainCtx, mainCancel := context.WithCancel(ctx)
	defer mainCancel()
	defer close(c.events)
	defer c.ticker.Stop()

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(c.config.PingTimeout))
	})

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		readCtx, readCancel := context.WithCancel(mainCtx)
		defer mainCancel()
		defer readCancel()
		defer c.conn.Close()

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
		defer c.conn.Close()

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
				// app shutdown
				return
			}
		}
	})

	wg.Wait()
	return nil
}
