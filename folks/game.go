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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ponyo877/folks-ui/entity"
	"github.com/ponyo877/folks-ui/websocket"
)

type Game struct {
	ws          *websocket.WebSocket
	schema      string
	host        string
	id          string
	x           int
	y           int
	now         time.Time
	messages    []*Message
	characters  map[string]*Character
	strokes     map[*Stroke]struct{}
	touchIDs    []ebiten.TouchID
	textField   *TextField
	messageArea *MessageArea
}

func NewGame(schema, host string) *Game {
	g := &Game{
		schema: schema,
		host:   host,
	}
	g.init()
	return g
}

func (g *Game) init() {
	g.now = time.Now()
	var err error
	if g.ws, err = websocket.NewWebSocket(g.schema, g.host, "/v1/socket"); err != nil {
		fmt.Printf("failed to connect to websocket: %v\n", err)
	}
	go g.ws.Receive(func(message *entity.Message) {
		id := message.Body().ID()
		switch message.MessageType() {
		case "enter":
			for _, user := range message.Body().Users() {
				g.characters[user.ID()] = NewCharacter(
					user.ID(),
					characterImage[user.ImgIdx()],
					user.X(),
					user.Y(),
					user.Dir(),
				)
			}
		case "move":
			g.characters[id] = NewCharacter(
				id,
				characterImage[message.Body().ImgIdx()],
				message.Body().X(),
				message.Body().Y(),
				message.Body().Dir(),
			)
		case "say":
			// ブラウザでの表示バグの暫定対処のために先頭にスペースを追加
			message, _ := NewMessage(id, " "+message.Body().Text(), message.CreatedAt())
			g.messageArea.AddMessage(message)
			g.messages = append(g.messages, message)
		case "leave":
			delete(g.characters, id)
		}
	})
	uuid, _ := uuid.NewRandom()
	g.id = uuid.String()
	g.strokes = map[*Stroke]struct{}{}

	w, h := characterImage[0].Bounds().Dx(), characterImage[0].Bounds().Dy()
	g.characters = map[string]*Character{}
	x, y, dir := rand.Intn(ScreenWidth-w), rand.Intn(ScreenHeight-h), entity.DirRight
	g.characters[g.id] = NewCharacter(
		g.id,
		characterImage[0],
		x,
		y,
		dir,
	)
	g.ws.Send(entity.NewMessage("enter", entity.NewEnterReqBody(g.id), g.now))
	g.ws.Send(entity.NewMessage("move", entity.NewMoveBody(g.id, x, y, dir), g.now))
}

func (g *Game) Exit() error {
	g.ws.Send(entity.NewMessage("leave", entity.NewLeaveBody(g.id), g.now))
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
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
		g.ws.Send(entity.NewMessage("move", entity.NewMoveBody(g.id, x, y, dir), g.now))
	}

	// message
	latestExpiredMessageIndex := -1
	for i, message := range g.messages {
		if message.IsExpired(g.now) {
			latestExpiredMessageIndex = i + 1
			continue
		}
	}
	if g.messageArea == nil {
		g.messageArea = NewMessageArea(MessageAreaPointX, 0)
	}
	if latestExpiredMessageIndex > 0 {
		// g.messageArea.AddMessage(g.messages[0:latestExpiredMessageIndex])
		g.messages = slices.Delete(g.messages, 0, latestExpiredMessageIndex)
	}

	// InputField
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		text := g.textField.Text()
		g.textField.Clear()
		g.ws.Send(entity.NewMessage("say", entity.NewSayBody(g.id, text), g.now))
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

	// Click & Touch
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := NewStroke(&MouseStrokeSource{})
		g.strokes[s] = struct{}{}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		s := NewStroke(&TouchStrokeSource{id})
		g.strokes[s] = struct{}{}
	}
	for s := range g.strokes {
		x, y := s.Position()
		dir := g.characters[g.id].Dir()
		g.ws.Send(entity.NewMessage("move", entity.NewMoveBody(g.id, x, y, dir), g.now))
		delete(g.strokes, s)
	}
	return nil
}

func (g *Game) Draw(Screen *ebiten.Image) {
	Screen.Fill(color.RGBA{0xeb, 0xeb, 0xeb, 0x01})
	g.drawGopher(Screen)
	g.drawTextField(Screen)
	g.drawMessageArea(Screen)
}

func (g *Game) drawGopher(screen *ebiten.Image) {
	characterKeys := make([]string, 0)
	for k := range g.characters {
		characterKeys = append(characterKeys, k)
	}
	sort.Strings(characterKeys)
	for _, k := range characterKeys {
		g.characters[k].Draw(screen, 0, 0, 1)
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

func (g *Game) drawMessageArea(screen *ebiten.Image) {
	g.messageArea.Draw(screen)
}
