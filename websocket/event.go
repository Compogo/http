package websocket

import "encoding/json"

type Event struct {
	Type      string          `json:"type,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp *Timestamp      `json:"timestamp,omitempty"`
}

func NewEvent(Type string, payload json.RawMessage) *Event {
	return &Event{Type: Type, Payload: payload}
}
