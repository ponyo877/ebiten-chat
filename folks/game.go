package folks

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
)

type Game struct {
	// The gopher's position
	x16 int
	y16 int

	now time.Time

	// Camera
	cameraX int
	cameraY int

	messages []*Message

	// Pipes
	pipeTileYs []int

	audioContext *audio.Context
	jumpPlayer   *audio.Player
	hitPlayer    *audio.Player

	textField *TextField
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
	g.cameraX = 0
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	g.x16, g.y16 = ebiten.CursorPosition()
	g.now = time.Now()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		message, _ := NewMessage("Hello, World!")
		g.messages = append(g.messages, message)
	}
	LatestExpiredMessageIndex := -1
	for i, message := range g.messages {
		if message.IsExpired(g.now) {
			LatestExpiredMessageIndex = i + 1
			continue
		}
	}
	if LatestExpiredMessageIndex > 0 {
		g.messages = slices.Delete(g.messages, 0, LatestExpiredMessageIndex)
	}

	if g.textField == nil {
		pX := 16
		pY := ScreenHeight - pX - textFieldHeight
		g.textField = NewTextField(image.Rect(pX, pY, ScreenWidth-pX, pY+textFieldHeight), false)
	}
	g.textField.Update()
	if g.textField.Contains(g.x16, g.y16) {
		g.textField.Focus()
		g.textField.SetSelectionStartByCursorPosition(g.x16, g.y16)
		ebiten.SetCursorShape(ebiten.CursorShapeText)
	} else {
		g.textField.Blur()
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	return nil
}

func (g *Game) Draw(Screen *ebiten.Image) {
	Screen.Fill(color.RGBA{0xeb, 0xeb, 0xeb, 0x01})
	g.drawGopher(Screen)
	g.drawTextField(Screen)

	ebitenutil.DebugPrint(Screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) drawGopher(Screen *ebiten.Image) {
	gopherX := g.x16 - g.cameraX
	gopherY := g.y16 - g.cameraY
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gopherX), float64(gopherY))
	op.Filter = ebiten.FilterLinear
	Screen.DrawImage(gopherImage, op)
	for _, message := range g.messages {
		speechBubble, _ := NewSpeechBubble(message, gopherX, gopherY)
		speechBubble.Draw(Screen, g.now)
	}
}

func (g *Game) drawTextField(screen *ebiten.Image) {
	g.textField.Draw(screen)
}
