package folks

import "time"

var (
	lifespan int64 = 3000
)

type Message struct {
	characterID string
	content     string
	createdAt   time.Time
}

// NewMessage creates a new Message
func NewMessage(characterID, content string) (*Message, error) {
	createdAt := time.Now()
	return &Message{
		characterID,
		content,
		createdAt,
	}, nil
}

// CharacterID returns the characterID of the Message
func (m *Message) CharacterID() string {
	return m.characterID
}

// Content returns the content of the Message
func (m *Message) Content() string {
	return m.content
}

// Size returns the size of the Message
func (m *Message) Size() float32 {
	return float32(len(m.content))
}

func (m *Message) ElapsedMilli(now time.Time) int64 {
	return now.UnixMilli() - m.createdAt.UnixMilli()
}

// IsExpired returns true if the Message is expired
func (m *Message) IsExpired(now time.Time) bool {
	return m.ElapsedMilli(now) > lifespan
}
