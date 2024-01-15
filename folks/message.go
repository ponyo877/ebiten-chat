package folks

import "time"

var (
	lifespan int64 = 10000000 // 10 seconds
)

type Message struct {
	content   string
	createdAt time.Time
}

// NewMessage creates a new Message
func NewMessage(content string) (*Message, error) {
	createdAt := time.Now()
	return &Message{
		content,
		createdAt,
	}, nil
}

// Content returns the content of the Message
func (m *Message) Content() string {
	return m.content
}

// Size returns the size of the Message
func (m *Message) Size() float32 {
	return float32(len(m.content))
}

func (m *Message) ElapsedMicro(now time.Time) int64 {
	return now.UnixMicro() - m.createdAt.UnixMicro()
}

// IsExpired returns true if the Message is expired
func (m *Message) IsExpired(now time.Time) bool {
	return m.ElapsedMicro(now) > lifespan
}
