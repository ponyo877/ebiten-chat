package folks

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
)

func floorDiv(x, y int) int {
	d := x / y
	if d*y == x || x >= 0 {
		return d
	}
	return d - 1
}

func floorMod(x, y int) int {
	return x - floorDiv(x, y)*y
}

type Game struct {
	mode Mode

	// The gopher's position
	x16  int
	y16  int
	vy16 int

	// Camera
	cameraX int
	cameraY int

	// Pipes
	pipeTileYs []int

	gameoverCount int

	touchIDs   []ebiten.TouchID
	gamepadIDs []ebiten.GamepadID

	audioContext *audio.Context
	jumpPlayer   *audio.Player
	hitPlayer    *audio.Player
}

func NewGame(crt bool) ebiten.Game {
	g := &Game{}
	g.init()
	if crt {
		return &GameWithCRTEffect{Game: g}
	}
	return g
}

func (g *Game) init() {
	g.x16 = 0
	g.y16 = 100 * 16
	g.cameraX = -240
	g.cameraY = 0
	g.pipeTileYs = make([]int, 256)
	for i := range g.pipeTileYs {
		g.pipeTileYs[i] = rand.Intn(6) + 2
	}

	if g.audioContext == nil {
		g.audioContext = audio.NewContext(48000)
	}

	jumpD, err := vorbis.DecodeWithoutResampling(bytes.NewReader(raudio.Jump_ogg))
	if err != nil {
		log.Fatal(err)
	}
	g.jumpPlayer, err = g.audioContext.NewPlayer(jumpD)
	if err != nil {
		log.Fatal(err)
	}

	jabD, err := wav.DecodeWithoutResampling(bytes.NewReader(raudio.Jab_wav))
	if err != nil {
		log.Fatal(err)
	}
	g.hitPlayer, err = g.audioContext.NewPlayer(jabD)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) isKeyJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	if len(g.touchIDs) > 0 {
		return true
	}
	g.gamepadIDs = ebiten.AppendGamepadIDs(g.gamepadIDs[:0])
	for _, g := range g.gamepadIDs {
		if ebiten.IsStandardGamepadLayoutAvailable(g) {
			if inpututil.IsStandardGamepadButtonJustPressed(g, ebiten.StandardGamepadButtonRightBottom) {
				return true
			}
			if inpututil.IsStandardGamepadButtonJustPressed(g, ebiten.StandardGamepadButtonRightRight) {
				return true
			}
		} else {
			// The button 0/1 might not be A/B buttons.
			if inpututil.IsGamepadButtonJustPressed(g, ebiten.GamepadButton0) {
				return true
			}
			if inpututil.IsGamepadButtonJustPressed(g, ebiten.GamepadButton1) {
				return true
			}
		}
	}
	return false
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if g.isKeyJustPressed() {
			g.mode = ModeGame
		}
	case ModeGame:
		g.x16 += 32
		g.cameraX += 2
		if g.isKeyJustPressed() {
			g.vy16 = -96
			if err := g.jumpPlayer.Rewind(); err != nil {
				return err
			}
			g.jumpPlayer.Play()
		}
		g.y16 += g.vy16

		// Gravity
		g.vy16 += 4
		if g.vy16 > 96 {
			g.vy16 = 96
		}

		if g.hit() {
			if err := g.hitPlayer.Rewind(); err != nil {
				return err
			}
			g.hitPlayer.Play()
			g.mode = ModeGameOver
			g.gameoverCount = 30
		}
	case ModeGameOver:
		if g.gameoverCount > 0 {
			g.gameoverCount--
		}
		if g.gameoverCount == 0 && g.isKeyJustPressed() {
			g.init()
			g.mode = ModeTitle
		}
	}
	return nil
}

func (g *Game) Draw(Screen *ebiten.Image) {
	Screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	g.drawTiles(Screen)
	if g.mode != ModeTitle {
		g.drawGopher(Screen)
	}

	var titleTexts string
	var texts string
	switch g.mode {
	case ModeTitle:
		titleTexts = "FLAPPY GOPHER"
		texts = "\n\n\n\n\n\nPRESS SPACE KEY\n\nOR A/B BUTTON\n\nOR TOUCH Screen"
	case ModeGameOver:
		texts = "\nGAME OVER!"
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(ScreenWidth/2, 3*titleFontSize)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = titleFontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(Screen, titleTexts, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   titleFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(ScreenWidth/2, 3*titleFontSize)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(Screen, texts, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize,
	}, op)

	if g.mode == ModeTitle {
		const msg = "Go Gopher by Renee French is\nlicenced under CC BY 3.0."

		op := &text.DrawOptions{}
		op.GeoM.Translate(ScreenWidth/2, ScreenHeight-smallFontSize/2)
		op.ColorScale.ScaleWithColor(color.White)
		op.LineSpacing = smallFontSize
		op.PrimaryAlign = text.AlignCenter
		op.SecondaryAlign = text.AlignEnd
		text.Draw(Screen, msg, &text.GoTextFace{
			Source: arcadeFaceSource,
			Size:   smallFontSize,
		}, op)
	}

	op = &text.DrawOptions{}
	op.GeoM.Translate(ScreenWidth, 0)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignEnd
	text.Draw(Screen, fmt.Sprintf("%04d", g.score()), &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize,
	}, op)

	ebitenutil.DebugPrint(Screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) pipeAt(tileX int) (tileY int, ok bool) {
	if (tileX - pipeStartOffsetX) <= 0 {
		return 0, false
	}
	if floorMod(tileX-pipeStartOffsetX, pipeIntervalX) != 0 {
		return 0, false
	}
	idx := floorDiv(tileX-pipeStartOffsetX, pipeIntervalX)
	return g.pipeTileYs[idx%len(g.pipeTileYs)], true
}

func (g *Game) score() int {
	x := floorDiv(g.x16, 16) / tileSize
	if (x - pipeStartOffsetX) <= 0 {
		return 0
	}
	return floorDiv(x-pipeStartOffsetX, pipeIntervalX)
}

func (g *Game) hit() bool {
	if g.mode != ModeGame {
		return false
	}
	const (
		gopherWidth  = 30
		gopherHeight = 60
	)
	w, h := gopherImage.Bounds().Dx(), gopherImage.Bounds().Dy()
	x0 := floorDiv(g.x16, 16) + (w-gopherWidth)/2
	y0 := floorDiv(g.y16, 16) + (h-gopherHeight)/2
	x1 := x0 + gopherWidth
	y1 := y0 + gopherHeight
	if y0 < -tileSize*4 {
		return true
	}
	if y1 >= ScreenHeight-tileSize {
		return true
	}
	xMin := floorDiv(x0-pipeWidth, tileSize)
	xMax := floorDiv(x0+gopherWidth, tileSize)
	for x := xMin; x <= xMax; x++ {
		y, ok := g.pipeAt(x)
		if !ok {
			continue
		}
		if x0 >= x*tileSize+pipeWidth {
			continue
		}
		if x1 < x*tileSize {
			continue
		}
		if y0 < y*tileSize {
			return true
		}
		if y1 >= (y+pipeGapY)*tileSize {
			return true
		}
	}
	return false
}

func (g *Game) drawTiles(Screen *ebiten.Image) {
	const (
		nx           = ScreenWidth / tileSize
		ny           = ScreenHeight / tileSize
		pipeTileSrcX = 128
		pipeTileSrcY = 192
	)

	op := &ebiten.DrawImageOptions{}
	for i := -2; i < nx+1; i++ {
		// ground
		op.GeoM.Reset()
		op.GeoM.Translate(float64(i*tileSize-floorMod(g.cameraX, tileSize)),
			float64((ny-1)*tileSize-floorMod(g.cameraY, tileSize)))
		Screen.DrawImage(tilesImage.SubImage(image.Rect(0, 0, tileSize, tileSize)).(*ebiten.Image), op)

		// pipe
		if tileY, ok := g.pipeAt(floorDiv(g.cameraX, tileSize) + i); ok {
			for j := 0; j < tileY; j++ {
				op.GeoM.Reset()
				op.GeoM.Scale(1, -1)
				op.GeoM.Translate(float64(i*tileSize-floorMod(g.cameraX, tileSize)),
					float64(j*tileSize-floorMod(g.cameraY, tileSize)))
				op.GeoM.Translate(0, tileSize)
				var r image.Rectangle
				if j == tileY-1 {
					r = image.Rect(pipeTileSrcX, pipeTileSrcY, pipeTileSrcX+tileSize*2, pipeTileSrcY+tileSize)
				} else {
					r = image.Rect(pipeTileSrcX, pipeTileSrcY+tileSize, pipeTileSrcX+tileSize*2, pipeTileSrcY+tileSize*2)
				}
				Screen.DrawImage(tilesImage.SubImage(r).(*ebiten.Image), op)
			}
			for j := tileY + pipeGapY; j < ScreenHeight/tileSize-1; j++ {
				op.GeoM.Reset()
				op.GeoM.Translate(float64(i*tileSize-floorMod(g.cameraX, tileSize)),
					float64(j*tileSize-floorMod(g.cameraY, tileSize)))
				var r image.Rectangle
				if j == tileY+pipeGapY {
					r = image.Rect(pipeTileSrcX, pipeTileSrcY, pipeTileSrcX+pipeWidth, pipeTileSrcY+tileSize)
				} else {
					r = image.Rect(pipeTileSrcX, pipeTileSrcY+tileSize, pipeTileSrcX+pipeWidth, pipeTileSrcY+tileSize+tileSize)
				}
				Screen.DrawImage(tilesImage.SubImage(r).(*ebiten.Image), op)
			}
		}
	}
}

func (g *Game) drawGopher(Screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := gopherImage.Bounds().Dx(), gopherImage.Bounds().Dy()
	op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	op.GeoM.Rotate(float64(g.vy16) / 96.0 * math.Pi / 6)
	op.GeoM.Translate(float64(w)/2.0, float64(h)/2.0)
	op.GeoM.Translate(float64(g.x16/16.0)-float64(g.cameraX), float64(g.y16/16.0)-float64(g.cameraY))
	op.Filter = ebiten.FilterLinear
	Screen.DrawImage(gopherImage, op)
}
