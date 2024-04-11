package drawable

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ponyo877/folks-ui/entity"
)

var maxMessagesCount = 50

type MessageArea struct {
	messages []*entity.ChatMessage
	x        float64
	y        float64
}

var (
	messageAreaHeight          float32     = LogFontSize*float32(maxMessagesCount) + 5
	messageAreaBackgroundColor color.Color = color.RGBA{0x22, 0x22, 0x22, 0x88}
)

func NewMessageArea(x, y int) *MessageArea {
	return &MessageArea{
		messages: []*entity.ChatMessage{},
		x:        float64(x),
		y:        float64(y),
	}
}

func (s *MessageArea) AddMessage(msg *entity.ChatMessage) {
	s.messages = append(s.messages, msg)
}

func (s *MessageArea) TruncateMessage() {
	s.messages = nil
}

func (s *MessageArea) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(s.x), float32(s.y), MessageAreaWidth, messageAreaHeight, messageAreaBackgroundColor, true)
	limit := min(len(s.messages), maxMessagesCount)
	for i, msg := range s.messages[len(s.messages)-limit:] {
		msgTxt := NewText(s.x, s.y+float64(i)*LogFontSize, LogFontSize, msg.Format(), color.Gray16{0xffff}, ArcadeFaceSource)
		msgTxt.Draw(screen, false)
	}
}
