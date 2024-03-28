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
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	d "github.com/ponyo877/folks-ui/drawable"
	"github.com/ponyo877/folks-ui/entity"
	"github.com/ponyo877/folks-ui/websocket"
)

type Game struct {
	schema       string
	host         string
	mode         Mode
	now          time.Time
	id           string
	x            int
	y            int
	dir          entity.Dir
	name         string
	imgid        int
	pushedDir    entity.Dir
	messages     []*entity.ChatMessage
	characters   map[string]*d.Character
	strokes      map[*d.Stroke]struct{}
	messageArea  *d.MessageArea
	messageField *d.TextField
	nameField    *d.TextField
	enterButton  *d.Button

	touchIDs []ebiten.TouchID
	ws       *websocket.WebSocket

	bluredX  int
	bluredY  int
	clickedX int
	clickedY int
}

const (
	NumOfImagesPerRow = 5
	Spacing           = 10
	cellSize          = 80 + Spacing
	paddingSelectY    = 50
	paddingNameY      = 175
	paddingEnterY     = 275
)

var (
	startNameX   = d.ScreenWidth/2 - (d.NameFieldWidth)/2
	startNameY   = d.ScreenHeight/2 - (d.NameFieldHeight)/2 - paddingNameY
	startSelectX = d.ScreenWidth/2 - (cellSize*NumOfImagesPerRow)/2
	startSelectY = d.ScreenHeight/2 - (cellSize*len(d.CharacterImage)/NumOfImagesPerRow)/2 + paddingSelectY
	startEnterX  = d.ScreenWidth / 2
	startEnterY  = d.ScreenHeight/2 + paddingEnterY
)

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
	w, h := d.CharacterImage[0].Bounds().Dx(), d.CharacterImage[0].Bounds().Dy()
	g.x, g.y, g.dir = rand.Intn(d.ScreenWidth-w), rand.Intn(d.ScreenHeight-h), entity.DirRight
	g.characters = map[string]*d.Character{}
	g.strokes = map[*d.Stroke]struct{}{}
	g.messageArea = d.NewMessageArea(d.MessageAreaPointX, 0)
	g.messageField = d.NewTextField(image.Rect(0, d.MessageFieldPointY, d.MessageAreaPointX, d.ScreenHeight))
	g.nameField = d.NewTextField(image.Rect(int(startNameX), int(startNameY), int(startNameX)+d.NameFieldWidth, int(startNameY)+d.NameFieldHeight))
	g.enterButton = d.NewButton("   ENTER   ", float64(startEnterX), float64(startEnterY), color.RGBA{0, 255, 0, 255})
	g.imgid = -1
	g.bluredX, g.bluredY, g.clickedX, g.clickedY = -1, -1, -1, -1

	var err error
	if g.ws, err = websocket.NewWebSocket(g.schema, g.host, "/v1/socket"); err != nil {
		fmt.Printf("failed to connect to websocket: %v\n", err)
	}
}

func (g *Game) Start() error {
	ebiten.SetWindowSize(d.ScreenWidth, d.ScreenHeight)
	ebiten.SetWindowTitle("Ebiten Chat")
	return ebiten.RunGame(g)
}

func (g *Game) Exit() error {
	g.ws.Send(entity.NewSocketMessage("leave", entity.NewLeaveBody(g.id), g.now))
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return d.ScreenWidth, d.ScreenHeight
}

func (g *Game) Update() error {
	g.x, g.y = ebiten.CursorPosition()
	g.now = time.Now()
	if g.mode == ModeTitle {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.mode = ModeChat
		}
		g.updateTitle()
		return nil
	}
	g.updateChat()
	return nil
}

func (g *Game) updateTitle() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.mode = ModeChat
	}
	// 名前フィールド
	g.name = g.nameField.Text()
	g.nameField.Update()
	if g.nameField.Contains(g.x, g.y) {
		g.nameField.Focus()
		g.nameField.SetSelectionStartByCursorPosition(g.x, g.y)
		ebiten.SetCursorShape(ebiten.CursorShapeText)
		return
	}
	g.nameField.Blur()
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)

	// キャラクター選択
	// マウスが当たっているセルの位置を計算
	if g.x-startSelectX > 0 && g.y-startSelectY > 0 {
		g.bluredX, g.bluredY = (g.x-startSelectX)/cellSize, (g.y-startSelectY)/cellSize
	}

	// マウスがクリックされたら、その位置を保存
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.clickedX, g.clickedY = g.bluredX, g.bluredY
		imgid := g.clickedX + g.clickedY*NumOfImagesPerRow
		if 0 <= imgid && imgid < len(d.CharacterImage) {
			g.imgid = g.clickedX + g.clickedY*NumOfImagesPerRow
		}
	}

	// 入室ボタン
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := d.NewStroke(&d.MouseStrokeSource{})
		g.strokes[s] = struct{}{}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		s := d.NewStroke(&d.TouchStrokeSource{ID: id})
		g.strokes[s] = struct{}{}
	}
	for s := range g.strokes {
		x, y := s.Position()
		if g.enterButton.Contains(x, y) {
			if g.name == "" {
				g.name = "名無し"
			}
			if g.imgid < 0 {
				g.imgid = rand.Intn(len(d.CharacterImage))
			}
			g.ws.Send(entity.NewSocketMessage("enter", entity.NewEnterReqBody(g.id), g.now))
			g.ws.Send(entity.NewSocketMessage("move", entity.NewMoveBody(g.id, g.x, g.y, g.name, g.imgid, g.dir), g.now))
			go g.ws.Receive(g.recieveMessage)
			g.mode = ModeChat
			return
		}
		delete(g.strokes, s)
	}
}

func (g *Game) updateChat() {
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
		g.ws.Send(entity.NewSocketMessage("move", entity.NewMoveBody(g.id, x, y, g.name, g.imgid, g.dir), g.now))
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

	// メッセージフィールド
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		text := g.messageField.Text()
		g.messageField.Clear()
		if strings.TrimSpace(text) != "" {
			g.ws.Send(entity.NewSocketMessage("say", entity.NewSayBody(g.id, text), g.now))
		}
	}
	g.messageField.Update()
	if g.messageField.Contains(g.x, g.y) {
		g.messageField.Focus()
		g.messageField.SetSelectionStartByCursorPosition(g.x, g.y)
		ebiten.SetCursorShape(ebiten.CursorShapeText)
		return
	}
	g.messageField.Blur()
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)

	// クリックと画面タッチ
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := d.NewStroke(&d.MouseStrokeSource{})
		g.strokes[s] = struct{}{}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		s := d.NewStroke(&d.TouchStrokeSource{ID: id})
		g.strokes[s] = struct{}{}
	}
	for s := range g.strokes {
		x, y := s.Position()
		g.ws.Send(entity.NewSocketMessage("move", entity.NewMoveBody(g.id, x, y, g.name, g.imgid, g.dir), g.now))
		delete(g.strokes, s)
	}
}

func (g *Game) Draw(Screen *ebiten.Image) {
	Screen.Fill(color.RGBA{0xeb, 0xeb, 0xeb, 0x01})
	if g.mode == ModeTitle {
		g.drawNameField(Screen)
		g.drawCharacterSelectArea(Screen)
		g.drawEnterButton(Screen)
		return
	}
	g.drawCharacters(Screen)
	g.drawSpeechBubble(Screen)
	g.drawMessageArea(Screen)
	g.drawMessageField(Screen)
}

func (g *Game) drawCharacters(screen *ebiten.Image) {
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
				speechBubble, _ := d.NewSpeechBubble(message, x, y)
				speechBubble.Draw(screen, g.now)
			}
		}
	}
}

func (g *Game) drawMessageArea(screen *ebiten.Image) {
	g.messageArea.Draw(screen)
}

func (g *Game) drawMessageField(screen *ebiten.Image) {
	g.messageField.Draw(screen)
}

func (g *Game) drawNameField(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("NAME: %v", g.name))
	g.nameField.Draw(screen)
}

func (g *Game) drawCharacterSelectArea(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nCHARACTERID: %v", g.imgid))
	// 格子を描画
	for i := 0; i < NumOfImagesPerRow; i++ {
		for j := 0; j < len(d.CharacterImage)/NumOfImagesPerRow; j++ {
			var clr color.Color = color.White
			if i == g.bluredX && j == g.bluredY {
				clr = color.RGBA{0, 0, 255, 255} // 青色でハイライト
			}
			if i == g.clickedX && j == g.clickedY {
				clr = color.RGBA{255, 0, 0, 255} // 赤色でハイライト
			}
			if clr != color.White {
				// g.imgid = i + j*NumOfImagesPerRow
				vector.DrawFilledRect(screen, float32(i*cellSize+startSelectX), float32(j*cellSize+startSelectY), cellSize, cellSize, clr, false)
			}
			vector.StrokeRect(screen, float32(i*cellSize+startSelectX), float32(j*cellSize+startSelectY), cellSize, cellSize, 1, clr, false)
		}
	}
	for i, img := range d.CharacterImage {
		w, h := img.Bounds().Dx(), img.Bounds().Dy()
		x := (w+Spacing)*(i%NumOfImagesPerRow) + startSelectX
		y := (h+Spacing)*(i/NumOfImagesPerRow) + startSelectY
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(img, opts)
	}
}

func (g *Game) drawEnterButton(screen *ebiten.Image) {
	g.enterButton.Draw(screen)
}

func (g *Game) recieveMessage(message *entity.SocketMessage) {
	id := message.Body().ID()
	switch message.MessageType() {
	case "enter":
		for _, user := range message.Body().Users() {
			g.characters[user.ID()] = d.NewCharacter(
				user.ID(),
				user.X(),
				user.Y(),
				user.Name(),
				user.ImgID(),
				user.Dir(),
			)
		}
	case "move":
		g.characters[id] = d.NewCharacter(
			id,
			message.Body().X(),
			message.Body().Y(),
			message.Body().Name(),
			message.Body().ImgID(),
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
