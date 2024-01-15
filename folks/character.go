package folks

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Dir int

const (
	DirLeft Dir = iota
	DirRight
)

type Character struct {
	leftImage  *ebiten.Image
	rightImage *ebiten.Image
	pX         int
	pY         int
	dir        Dir
}

func NewCharacter(leftImage, rightImage *ebiten.Image) *Character {
	return &Character{
		leftImage:  leftImage,
		rightImage: rightImage,
		pX:         0,
		pY:         0,
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
	return c.leftImage.At(x-c.pX, y-c.pY).(color.RGBA).A > 0
}

// MoveBy moves the sprite by (x, y).
func (c *Character) MoveBy(x, y int) {
	w, h := c.leftImage.Bounds().Dx(), c.leftImage.Bounds().Dy()

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

func (c *Character) dirImage(dir Dir) *ebiten.Image {
	if dir == DirRight {
		return c.rightImage
	}
	return c.leftImage
}

// Draw draws the sprite.
func (c *Character) Draw(screen *ebiten.Image, dx, dy int, dir Dir, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.pX+dx), float64(c.pY+dy))
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(c.dirImage(dir), op)
	screen.DrawImage(c.dirImage(dir), op)
}
