package websocket

import "errors"

var (
	// MessageChanFullError возникает, когда буфер исходящих сообщений клиента заполнен.
	MessageChanFullError = errors.New("message channel is full")

	// MessageChanClosedError возникает при попытке записи в закрытый канал.
	MessageChanClosedError = errors.New("message channel is closed")
)
