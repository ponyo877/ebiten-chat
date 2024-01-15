package folks

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Character struct {
	image *ebiten.Image
	pX    int
	pY    int
}

func NewCharacter(image *ebiten.Image) *Character {
	return &Character{
		image: image,
		pX:    0,
		pY:    0,
	}
}

func (c *Character) Point() (int, int) {
	return c.pX, c.pY
}

func (c *Character) Move(x, y int) {
	c.pX = x
	c.pY = y
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (c *Character) In(x, y int) bool {
	// Check the actual color (alpha) value at the specified position
	// so that the result of In becomes natural to userc.
	//
	// Note that this is not a good manner to use At for logic
	// since color from At might include some errors on some machinec.
	// As this is not so important logic, it's ok to use it so far.
	return c.image.At(x-c.pX, y-c.pY).(color.RGBA).A > 0
}

// MoveBy moves the sprite by (x, y).
func (c *Character) MoveBy(x, y int) {
	w, h := c.image.Bounds().Dx(), c.image.Bounds().Dy()

	c.pX += x
	c.pY += y
	if c.pX < 0 {
		c.pX = 0
	}
	if c.pX > ScreenWidth-w {
		c.pX = ScreenWidth - w
	}
	if c.pY < 0 {
		c.pY = 0
	}
	if c.pY > ScreenHeight-h {
		c.pY = ScreenHeight - h
	}
}

// Draw draws the sprite.
func (c *Character) Draw(screen *ebiten.Image, dx, dy int, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.pX+dx), float64(c.pY+dy))
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(c.image, op)
	screen.DrawImage(c.image, op)
}
