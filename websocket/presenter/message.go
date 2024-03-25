package presenter

import (
	"time"

	"github.com/ponyo877/folks-ui/entity"
)

// MessagePresenter
type MessagePresenter struct {
	MessageType string         `json:"messageType"` // say, move, enter, leave
	Body        *BodyPresenter `json:"body"`
	CreatedAt   time.Time      `json:"createdAt"`
}

func NewMessagePresenter(message *entity.SocketMessage) MessagePresenter {
	return MessagePresenter{
		MessageType: message.MessageType(),
		Body:        NewBodyPresenter(message.MessageType(), message.Body()),
		CreatedAt:   message.CreatedAt(),
	}
}

// Unmarshal
func (m MessagePresenter) Unmarshal() *entity.SocketMessage {
	return entity.NewSocketMessage(m.MessageType, m.Body.Unmarshal(m.MessageType), m.CreatedAt)
}

// MarshalMessage
func MarshalMessage(message *entity.SocketMessage) MessagePresenter {
	return NewMessagePresenter(message)
}
