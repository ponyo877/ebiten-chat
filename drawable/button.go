package drawable

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	text      string
	x         float64
	y         float64
	clr       color.Color
	buttonTxt *Text
}

func NewButton(text string, x, y float64, clr color.Color) *Button {
	buttonTxt := NewText(x, y, fontSize, text, color.Black, arcadeFaceSource)
	return &Button{
		text,
		x,
		y,
		clr,
		buttonTxt,
	}
}

func (b *Button) Contains(x, y int) bool {
	w, h := b.buttonTxt.Bounds()
	return b.x-w/2 <= float64(x) && float64(x) <= b.x+w/2 && b.y <= float64(y) && float64(y) <= b.y+h
}

func (b *Button) Draw(screen *ebiten.Image) {
	w, h := b.buttonTxt.Bounds()

	vector.DrawFilledRect(screen, float32(b.x-w/2), float32(b.y), float32(w), float32(h), b.clr, true)
	b.buttonTxt.Draw(screen, true)
}
