package folks

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"slices"
	"sort"
	"strings"
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
	dir         entity.Dir
	pushedDir   entity.Dir
	now         time.Time
	messages    []*entity.ChatMessage
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
	go g.ws.Receive(func(message *entity.SocketMessage) {
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
			message, _ := entity.NewChatMessage(id, message.Body().Text(), message.CreatedAt())
			if message != nil {
				g.messageArea.AddMessage(message)
				g.messages = append(g.messages, message)
			}
		case "leave":
			delete(g.characters, id)
		}
	})
	uuid, _ := uuid.NewRandom()
	g.id = uuid.String()
	g.strokes = map[*Stroke]struct{}{}

	w, h := characterImage[0].Bounds().Dx(), characterImage[0].Bounds().Dy()
	g.characters = map[string]*Character{}
	g.x, g.y, g.dir = rand.Intn(ScreenWidth-w), rand.Intn(ScreenHeight-h), entity.DirRight
	g.messageArea = NewMessageArea(MessageAreaPointX, 0)
	pX := TextFieldPadding
	pY := ScreenHeight - TextFieldPadding - TextFieldHeight
	g.textField = NewTextField(image.Rect(pX, pY, ScreenWidth-pX, pY+TextFieldHeight), false)

	g.ws.Send(entity.NewSocketMessage("enter", entity.NewEnterReqBody(g.id), g.now))
	g.ws.Send(entity.NewSocketMessage("move", entity.NewMoveBody(g.id, g.x, g.y, g.dir), g.now))
}

func (g *Game) Exit() error {
	g.ws.Send(entity.NewSocketMessage("leave", entity.NewLeaveBody(g.id), g.now))
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	g.x, g.y = ebiten.CursorPosition()
	g.now = time.Now()

	// Character Direction
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.pushedDir = entity.DirRight
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.pushedDir = entity.DirLeft
	}
	if _, ok := g.characters[g.id]; ok && g.dir != g.pushedDir {
		x, y := g.characters[g.id].Point()
		g.dir = g.pushedDir
		g.ws.Send(entity.NewSocketMessage("move", entity.NewMoveBody(g.id, x, y, g.dir), g.now))
	}

	// Meessage
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

	// Input Field
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		text := g.textField.Text()
		g.textField.Clear()
		if strings.TrimSpace(text) != "" {
			g.ws.Send(entity.NewSocketMessage("say", entity.NewSayBody(g.id, text), g.now))
		}
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
		g.ws.Send(entity.NewSocketMessage("move", entity.NewMoveBody(g.id, x, y, g.dir), g.now))
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
