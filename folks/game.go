package folks

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"slices"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ponyo877/folks-ui/entity"
	"github.com/ponyo877/folks-ui/websocket"
)

type Game struct {
	wss        *websocket.WebSocket
	id         string
	x          int
	y          int
	now        time.Time
	messages   []*Message
	characters map[string]*Character
	strokes    map[*Stroke]struct{}
	touchIDs   []ebiten.TouchID
	textField  *TextField
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
	var err error
	if g.wss, err = websocket.NewWebSocket("folks-chat.com", "/v1/socket"); err != nil {
		fmt.Printf("failed to connect to websocket: %v\n", err)
	}
	go g.wss.Receive(func(message *entity.Message) {
		id := message.Body().ID()
		switch message.MessageType() {
		case "move":
			g.characters[id] = NewCharacter(
				id,
				gopherLeftImage,
				gopherRightImage,
				message.Body().X(),
				message.Body().Y(),
				message.Body().Dir(),
			)
		case "say":
			message, _ := NewMessage(id, message.Body().Text())
			g.messages = append(g.messages, message)
		}
	})
	uuid, _ := uuid.NewRandom()
	g.id = uuid.String()
	g.strokes = map[*Stroke]struct{}{}

	w, h := gopherLeftImage.Bounds().Dx(), gopherLeftImage.Bounds().Dy()
	g.characters = map[string]*Character{}
	x, y, dir := rand.Intn(ScreenWidth-w), rand.Intn(ScreenHeight-h), entity.DirRight
	g.characters[g.id] = NewCharacter(
		g.id,
		gopherLeftImage,
		gopherRightImage,
		x,
		y,
		dir,
	)
	g.wss.Send(entity.NewMessage("move", entity.NewMoveBody(g.id, x, y, dir)))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	g.x, g.y = ebiten.CursorPosition()
	g.now = time.Now()

	// character direction
	dir := entity.DirUnknown
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dir = entity.DirRight
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dir = entity.DirLeft
	}
	if dir != entity.DirUnknown {
		x, y := g.characters[g.id].Point()
		g.wss.Send(entity.NewMessage("move", entity.NewMoveBody(g.id, x, y, dir)))
	}

	// message
	latestExpiredMessageIndex := -1
	for i, message := range g.messages {
		if message.IsExpired(g.now) {
			latestExpiredMessageIndex = i + 1
			continue
		}
	}
	if latestExpiredMessageIndex > 0 {
		g.messages = slices.Delete(g.messages, 0, latestExpiredMessageIndex)
	}

	// InputField
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		text := g.textField.Text()
		g.textField.Clear()
		g.wss.Send(entity.NewMessage("say", entity.NewSayBody(g.id, text)))
	}
	if g.textField == nil {
		pX := TextFieldPadding
		pY := ScreenHeight - TextFieldPadding - TextFieldHeight
		g.textField = NewTextField(image.Rect(pX, pY, ScreenWidth-pX, pY+TextFieldHeight), false)
	}
	g.textField.Update()
	if g.textField.Contains(g.x, g.y) {
		g.textField.Focus()
		g.textField.SetSelectionStartByCursorPosition(g.x, g.y)
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
			x, y := s.Position()
			dir := g.characters[g.id].Dir()
			g.wss.Send(entity.NewMessage("move", entity.NewMoveBody(g.id, x, y, dir)))
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
	draggingCharacters := map[*Character]struct{}{}
	for s := range g.strokes {
		if character := s.DraggingObject().(*Character); character != nil {
			draggingCharacters[character] = struct{}{}
		}
	}

	keys := make([]string, 0)
	for k := range g.characters {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		c := g.characters[k]
		if _, ok := draggingCharacters[c]; ok {
			continue
		}
		c.Draw(screen, 0, 0, 1)
	}
	for s := range g.strokes {
		dx, dy := s.PositionDiff()
		if c := s.DraggingObject().(*Character); c != nil {
			c.Draw(screen, dx, dy, 0.5)
		}
	}

	// SpeechBubble
	for _, character := range g.characters {
		x, y := character.Point()
		for _, message := range g.messages {
			if character.IsMine(message.CharacterID()) {
				speechBubble, _ := NewSpeechBubble(message, x, y)
				speechBubble.Draw(screen, g.now)
			}
		}
	}
}

func (g *Game) drawTextField(screen *ebiten.Image) {
	g.textField.Draw(screen)
}

func (g *Game) characterAt(x, y int) *Character {
	// As the characters are ordered from back to front,
	// search the clicked/touched character in reverse order.
	for _, c := range g.characters {
		if c.In(x, y) && c.IsMine(g.id) {
			return c
		}
	}
	return nil
}

func (g *Game) updateStroke(stroke *Stroke) {
	stroke.Update()
	if !stroke.IsReleased() {
		return
	}

	c := stroke.DraggingObject().(*Character)
	if c == nil {
		return
	}

	c.MoveBy(stroke.PositionDiff())
	stroke.SetDraggingObject(nil)
}