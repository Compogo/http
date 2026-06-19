package websocket

import "encoding/json"

// Event представляет событие, передаваемое по WebSocket.
// Содержит тип, полезную нагрузку (payload) и временную метку.
//
// Пример:
//
//	event := &Event{
//	    Type: "chat_message",
//	    Payload: json.RawMessage(`{"text":"hello","user":"john"}`),
//	}
type Event struct {
	Type      string          `json:"type,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp *Timestamp      `json:"timestamp,omitempty"`
}

// NewEvent создаёт новое событие с указанным типом и payload'ом.
func NewEvent(Type string, payload json.RawMessage) *Event {
	return &Event{Type: Type, Payload: payload}
}
