package websocket

import "errors"

var (
	MessageChanFullError   = errors.New("message channel is full")
	MessageChanClosedError = errors.New("message channel is closed")
)
