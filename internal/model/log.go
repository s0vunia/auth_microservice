package model

import "time"

// Log represents a log entity with ID, Message, and CreatedAt fields
type Log struct {
	ID        int64
	Message   string
	CreatedAt time.Time
}

// LogCreate represents a log create entity with Message field
type LogCreate struct {
	Message string
}
