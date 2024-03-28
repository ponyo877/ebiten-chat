package game

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
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
		g.updateNameField()
		g.updateCharacterSelect()
		g.updateEnterButton()
		return nil
	}
	g.updateCharacterDir()
	g.updateChatMsg()
	if skip := g.updateMsgField(); skip {
		return nil
	}
	g.updateMove()
	return nil
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
