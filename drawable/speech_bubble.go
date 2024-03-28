package drawable

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ponyo877/folks-ui/entity"
)

type SpeechBubble struct {
	message   *entity.ChatMessage
	x         float64
	y         float64
	createdAt int64
}

var (
	speechBubbleAltitudeUnit float64 = 0.03
)

func NewSpeechBubble(message *entity.ChatMessage, x, y int) (*SpeechBubble, error) {
	createdAt := time.Now().Unix()
	return &SpeechBubble{
		message:   message,
		x:         float64(x),
		y:         float64(y),
		createdAt: createdAt,
	}, nil
}

func (s *SpeechBubble) Draw(screen *ebiten.Image, now time.Time) {
	ch := float64(CharacterImage[0].Bounds().Dy())
	bh := s.y - ch - s.altitude(now)
	bubbleTxt := NewText(float64(s.x), bh, fontSize, s.message.Content(), color.Black, arcadeFaceSource)
	w, h := bubbleTxt.Bounds()

	vector.DrawFilledRect(screen, float32(s.x-w/2), float32(bh), float32(w), float32(h), color.White, true)
	bubbleTxt.Draw(screen, true)
}

func (s *SpeechBubble) altitude(now time.Time) float64 {
	return float64(s.message.ElapsedMilli(now)) * speechBubbleAltitudeUnit
}
