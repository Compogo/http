package websocket

import (
	"time"
)

// Timestamp представляет временную метку в формате Unix.
// Используется для отметки времени создания событий.
type Timestamp struct {
	Seconds int64 `json:"seconds,omitempty"`
	Nano    int32 `json:"nano,omitempty"`
}

// NewTimestamp создаёт временную метку для текущего момента времени (UTC).
func NewTimestamp() *Timestamp {
	now := time.Now().UTC()

	return &Timestamp{
		Seconds: now.Unix(),
		Nano:    int32(now.Nanosecond()),
	}
}
