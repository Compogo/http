package websocket

import "encoding/json"

type Event struct {
	Type      string          `json:"type,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp *Timestamp      `json:"timestamp,omitempty"`
}

type Timestamp struct {
	Seconds int64 `json:"seconds,omitempty"`
	Nanos   int64 `json:"nanos,omitempty"`
}
