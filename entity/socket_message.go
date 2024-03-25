package entity

import (
	"time"
)

type SocketMessage struct {
	messageType string
	body        *Body
	createdAt   time.Time
}

// NewSocketMessage
func NewSocketMessage(messageType string, body *Body, createdAt time.Time) *SocketMessage {
	return &SocketMessage{
		messageType: messageType,
		body:        body,
		createdAt:   createdAt,
	}
}

// MessageType
func (m *SocketMessage) MessageType() string {
	return m.messageType
}

// Body
func (m *SocketMessage) Body() *Body {
	return m.body
}

// CreatedAt
func (m *SocketMessage) CreatedAt() time.Time {
	return m.createdAt
}
