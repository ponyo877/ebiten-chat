package drawable

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Stroke struct {
	x int
	y int
}

func NewStroke(source StrokeSource) *Stroke {
	x, y := source.Position()
	return &Stroke{x, y}
}

func (s *Stroke) Position() (int, int) {
	return s.x, s.y
}

type StrokeSource interface {
	Position() (int, int)
}

type MouseStrokeSource struct{}

func (m *MouseStrokeSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

type TouchStrokeSource struct {
	ID ebiten.TouchID
}

func (t *TouchStrokeSource) Position() (int, int) {
	return ebiten.TouchPosition(t.ID)
}
