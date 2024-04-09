package drawable

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	text     string
	x        float64
	y        float64
	fontSize float64
	w        float64
	h        float64
}

func NewButton(text string, x, y float64, fontSize float64) *Button {
	return &Button{
		text:     text,
		x:        x,
		y:        y,
		fontSize: fontSize,
	}
}

func (b *Button) Bounds() (float64, float64) {
	buttonTxt := NewText(b.x-float64(b.w)/2, b.y, b.fontSize, b.text, color.Black, arcadeFaceSource)
	return buttonTxt.Bounds()
}

func (b *Button) Contains(x, y int) bool {
	return b.x-b.w/2 <= float64(x) && float64(x) <= b.x+b.w/2 && b.y <= float64(y) && float64(y) <= b.y+b.h
}

func (b *Button) SetWH(w, h float64) {
	b.w, b.h = w, h
}

func (b *Button) FixDraw(screen *ebiten.Image, clr color.Color) {
	vector.DrawFilledRect(screen, float32(b.x-float64(b.w)/2), float32(b.y), float32(b.w), float32(b.h), clr, true)
	buttonTxt := NewText(b.x-float64(b.w)/2, b.y, b.fontSize, b.text, color.Black, arcadeFaceSource)
	buttonTxt.Draw(screen, false)
}
