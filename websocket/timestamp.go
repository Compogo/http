package websocket

import (
	"time"
)

type Timestamp struct {
	Seconds int64 `json:"seconds,omitempty"`
	Nano    int32 `json:"nano,omitempty"`
}

func NewTimestamp() *Timestamp {
	now := time.Now().UTC()

	return &Timestamp{
		Seconds: now.Unix(),
		Nano:    int32(now.Nanosecond()),
	}
}
