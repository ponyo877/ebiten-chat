package folks

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ponyo877/folks-ui/entity"
)

type Character struct {
	id         string
	leftImage  *ebiten.Image
	rightImage *ebiten.Image
	x          int
	y          int
	dir        entity.Dir
}

func NewCharacter(id string, leftImage, rightImage *ebiten.Image, x, y int, dir entity.Dir) *Character {
	return &Character{
		id:         id,
		leftImage:  leftImage,
		rightImage: rightImage,
		x:          x,
		y:          y,
		dir:        dir,
	}
}

func (c *Character) Point() (int, int) {
	return c.x, c.y
}

func (c *Character) Dir() entity.Dir {
	return c.dir
}

func (c *Character) Move(x, y int) {
	c.x = x
	c.y = y
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (c *Character) In(x, y int) bool {
	// Check the actual color (alpha) value at the specified position
	// so that the result of In becomes natural to userc.
	//
	// Note that this is not a good manner to use At for logic
	// since color from At might include some errors on some machinec.
	// As this is not so important logic, it's ok to use it so far.
	return c.leftImage.At(x-c.x, y-c.y).(color.RGBA).A > 0
}

// MoveBy moves the sprite by (x, y).
func (c *Character) MoveBy(x, y int) {
	w, h := c.leftImage.Bounds().Dx(), c.leftImage.Bounds().Dy()

	c.x += x
	c.y += y
	if c.x < 0 {
		c.x = 0
	}
	if c.x > ScreenWidth-w {
		c.x = ScreenWidth - w
	}
	if c.y < 0 {
		c.y = 0
	}
	if c.y > ScreenHeight-h {
		c.y = ScreenHeight - h
	}
}

func (c *Character) IsMine(myID string) bool {
	return c.id == myID
}

func (c *Character) dirImage() *ebiten.Image {
	if c.dir == entity.DirLeft {
		return c.leftImage
	}
	return c.rightImage
}

// Draw draws the sprite.
func (c *Character) Draw(screen *ebiten.Image, dx, dy int, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(float64(c.x+dx), float64(c.y+dy))
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(c.dirImage(), op)
}
