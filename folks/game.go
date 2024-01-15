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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	// The gopher's position
	x16 int
	y16 int

	now time.Time

	// Camera
	cameraX int
	cameraY int

	dir Dir

	messages []*Message

	audioContext *audio.Context
	jumpPlayer   *audio.Player
	hitPlayer    *audio.Player

	characters []*Character
	strokes    map[*Stroke]struct{}
	touchIDs   []ebiten.TouchID

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
	g.dir = DirRight

	w, h := gopherLeftImage.Bounds().Dx(), gopherLeftImage.Bounds().Dy()
	g.characters = append(g.characters, &Character{
		leftImage:  gopherLeftImage,
		rightImage: gopherRightImage,
		pX:         rand.Intn(ScreenWidth - w),
		pY:         rand.Intn(ScreenHeight - h),
	})
	g.strokes = map[*Stroke]struct{}{}

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
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		text := g.textField.Text()
		g.textField.Clear()
		message, _ := NewMessage(text)
		g.messages = append(g.messages, message)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.dir = DirRight
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.dir = DirLeft
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

	// InputField
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

	// Drug & Drop
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := NewStroke(&MouseStrokeSource{})
		s.SetDraggingObject(g.characterAt(s.Position()))
		g.strokes[s] = struct{}{}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		s := NewStroke(&TouchStrokeSource{id})
		s.SetDraggingObject(g.characterAt(s.Position()))
		g.strokes[s] = struct{}{}
	}
	for s := range g.strokes {
		g.updateStroke(s)
		if s.IsReleased() {
			delete(g.strokes, s)
		}
	}
	return nil
}

func (g *Game) Draw(Screen *ebiten.Image) {
	Screen.Fill(color.RGBA{0xeb, 0xeb, 0xeb, 0x01})
	g.drawGopher(Screen)
	g.drawTextField(Screen)

	ebitenutil.DebugPrint(Screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) drawGopher(screen *ebiten.Image) {
	// gopherX := g.x16 - g.cameraX
	// gopherY := g.y16 - g.cameraY

	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(float64(gopherX), float64(gopherY))
	// op.Filter = ebiten.FilterLinear
	// gopherImage := gopherRightImage
	// if g.dir == DirLeft {
	// 	gopherImage = gopherLeftImage
	// }
	// Screen.DrawImage(gopherImage, op)
	// for _, message := range g.messages {
	// 	speechBubble, _ := NewSpeechBubble(message, gopherX, gopherY)
	// 	speechBubble.Draw(Screen, g.now)
	// }
	draggingCharacters := map[*Character]struct{}{}
	for s := range g.strokes {
		if character := s.DraggingObject().(*Character); character != nil {
			draggingCharacters[character] = struct{}{}
		}
	}

	for _, s := range g.characters {
		if _, ok := draggingCharacters[s]; ok {
			continue
		}
		s.Draw(screen, 0, 0, g.dir, 1)
	}
	for s := range g.strokes {
		dx, dy := s.PositionDiff()
		if character := s.DraggingObject().(*Character); character != nil {
			character.Draw(screen, dx, dy, g.dir, 0.5)
		}
	}

	// SpeechBubble
	for _, character := range g.characters {
		gopherX, gopherY := character.Point()
		for _, message := range g.messages {
			speechBubble, _ := NewSpeechBubble(message, gopherX, gopherY)
			speechBubble.Draw(screen, g.now)
		}
	}

}

func (g *Game) drawTextField(screen *ebiten.Image) {
	g.textField.Draw(screen)
}

func (g *Game) characterAt(x, y int) *Character {
	// As the characters are ordered from back to front,
	// search the clicked/touched character in reverse order.
	for i := len(g.characters) - 1; i >= 0; i-- {
		s := g.characters[i]
		if s.In(x, y) {
			return s
		}
	}
	return nil
}

func (g *Game) updateStroke(stroke *Stroke) {
	stroke.Update()
	if !stroke.IsReleased() {
		return
	}

	s := stroke.DraggingObject().(*Character)
	if s == nil {
		return
	}

	s.MoveBy(stroke.PositionDiff())

	index := -1
	for i, ss := range g.characters {
		if ss == s {
			index = i
			break
		}
	}

	// Move the dragged character to the front.
	g.characters = append(g.characters[:index], g.characters[index+1:]...)
	g.characters = append(g.characters, s)

	stroke.SetDraggingObject(nil)
}
