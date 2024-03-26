package game

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
	"github.com/ponyo877/folks-ui/drawable"
	"github.com/ponyo877/folks-ui/entity"
	"github.com/ponyo877/folks-ui/websocket"
)

type Game struct {
	schema      string
	host        string
	now         time.Time
	id          string
	x           int
	y           int
	dir         entity.Dir
	pushedDir   entity.Dir
	messages    []*entity.ChatMessage
	characters  map[string]*drawable.Character
	strokes     map[*drawable.Stroke]struct{}
	messageArea *drawable.MessageArea
	textField   *drawable.TextField
	touchIDs    []ebiten.TouchID
	ws          *websocket.WebSocket
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
	uuid, _ := uuid.NewRandom()
	g.id = uuid.String()
	w, h := drawable.CharacterImage[0].Bounds().Dx(), drawable.CharacterImage[0].Bounds().Dy()
	g.x, g.y, g.dir = rand.Intn(drawable.ScreenWidth-w), rand.Intn(drawable.ScreenHeight-h), entity.DirRight
	g.characters = map[string]*drawable.Character{}
	g.strokes = map[*drawable.Stroke]struct{}{}
	g.messageArea = drawable.NewMessageArea(drawable.MessageAreaPointX, 0)
	g.textField = drawable.NewTextField(image.Rect(0, drawable.TextFieldPointY, drawable.MessageAreaPointX, drawable.ScreenHeight))

	var err error
	if g.ws, err = websocket.NewWebSocket(g.schema, g.host, "/v1/socket"); err != nil {
		fmt.Printf("failed to connect to websocket: %v\n", err)
	}
	go g.ws.Receive(g.recieveMessage)
	g.ws.Send(entity.NewSocketMessage("enter", entity.NewEnterReqBody(g.id), g.now))
	g.ws.Send(entity.NewSocketMessage("move", entity.NewMoveBody(g.id, g.x, g.y, g.dir), g.now))
}

func (g *Game) Start() error {
	ebiten.SetWindowSize(drawable.ScreenWidth, drawable.ScreenHeight)
	ebiten.SetWindowTitle("Ebiten Chat")
	return ebiten.RunGame(g)
}

func (g *Game) Exit() error {
	g.ws.Send(entity.NewSocketMessage("leave", entity.NewLeaveBody(g.id), g.now))
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return drawable.ScreenWidth, drawable.ScreenHeight
}

func (g *Game) Update() error {
	g.x, g.y = ebiten.CursorPosition()
	g.now = time.Now()

	// キャラの向き
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

	// チャットメッセージ
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

	// テキストフィールド
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
		return nil
	}
	g.textField.Blur()
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)

	// クリックと画面タッチ
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := drawable.NewStroke(&drawable.MouseStrokeSource{})
		g.strokes[s] = struct{}{}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		s := drawable.NewStroke(&drawable.TouchStrokeSource{ID: id})
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
	g.drawSpeechBubble(Screen)
	g.drawMessageArea(Screen)
	g.drawTextField(Screen)
}

func (g *Game) drawGopher(screen *ebiten.Image) {
	characterKeys := make([]string, 0)
	for k := range g.characters {
		characterKeys = append(characterKeys, k)
	}
	sort.Strings(characterKeys)
	for _, k := range characterKeys {
		g.characters[k].Draw(screen)
	}
}

func (g *Game) drawSpeechBubble(screen *ebiten.Image) {
	for _, character := range g.characters {
		x, y := character.Point()
		for _, message := range g.messages {
			if character.IsMine(message.CharacterID()) {
				speechBubble, _ := drawable.NewSpeechBubble(message, x, y)
				speechBubble.Draw(screen, g.now)
			}
		}
	}
}

func (g *Game) drawMessageArea(screen *ebiten.Image) {
	g.messageArea.Draw(screen)
}

func (g *Game) drawTextField(screen *ebiten.Image) {
	g.textField.Draw(screen)
}

func (g *Game) recieveMessage(message *entity.SocketMessage) {
	id := message.Body().ID()
	switch message.MessageType() {
	case "enter":
		for _, user := range message.Body().Users() {
			g.characters[user.ID()] = drawable.NewCharacter(
				user.ID(),
				drawable.CharacterImage[user.ImgIdx()],
				user.X(),
				user.Y(),
				user.Dir(),
			)
		}
	case "move":
		g.characters[id] = drawable.NewCharacter(
			id,
			drawable.CharacterImage[message.Body().ImgIdx()],
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
}
