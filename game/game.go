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
	roomID       string
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
	roomButtons  []*d.Button
	roomText     *d.Text
	undoText     *d.Text
	touchIDs     []ebiten.TouchID
	ws           *websocket.WebSocket

	bluredCharacterX  int
	bluredCharacterY  int
	clickedCharacterX int
	clickedCharacterY int
	bluredRoom        int
}

const (
	NumOfImagesPerRow = 5
	Spacing           = 10
	cellSize          = 80 + Spacing
	paddingSelectY    = 50
	paddingNameY      = 175
	paddingRoomY      = 20
)

var (
	startNameX   = d.ScreenWidth/4 - d.NameFieldWidth/2
	startNameY   = d.ScreenHeight/4 - d.NameFieldHeight/2
	startSelectX = int(d.ScreenWidth/4 - (cellSize*NumOfImagesPerRow)/2)
	startSelectY = int(startNameY + paddingSelectY)
	startRoomX   = int(d.ScreenWidth * 3 / 4)
	startRoomY   = int(startNameY)
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
	uuid, _ := uuid.NewRandom()
	g.id = uuid.String()
	w, h := d.CharacterImages[0].Bounds().Dx(), d.CharacterImages[0].Bounds().Dy()
	g.x, g.y, g.dir = rand.Intn(d.ScreenWidth-w), rand.Intn(d.ScreenHeight-h), entity.DirRight
	g.characters = map[string]*d.Character{}
	g.strokes = map[*d.Stroke]struct{}{}
	g.messageArea = d.NewMessageArea(d.MessageAreaPointX, 0)
	g.messageField = d.NewTextField(image.Rect(0, d.MessageFieldPointY, d.MessageAreaPointX, d.ScreenHeight))
	g.nameField = d.NewTextField(image.Rect(int(startNameX), int(startNameY), int(startNameX)+d.NameFieldWidth, int(startNameY)+d.NameFieldHeight))
	g.roomButtons = make([]*d.Button, len(d.RoomNameList))
	for i := 0; i < len(d.RoomNameList); i++ {
		roomText := fmt.Sprintf("  Room#%d %s  ", i+1, d.RoomNameList[i])
		largeFontSize := d.LargeFontSize
		g.roomButtons[i] = d.NewButton(roomText, float64(startRoomX), float64(startRoomY)+float64(i)*(largeFontSize+paddingRoomY), largeFontSize)
	}
	g.imgid = -1
	g.bluredCharacterX, g.bluredCharacterY, g.clickedCharacterX, g.clickedCharacterY = -1, -1, -1, -1
}

func (g *Game) Start() error {
	ebiten.SetWindowSize(d.ScreenWidth, d.ScreenHeight)
	ebiten.SetWindowTitle("Ebiten Chat")
	return ebiten.RunGame(g)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return d.ScreenWidth, d.ScreenHeight
}

func (g *Game) Update() error {
	g.x, g.y = ebiten.CursorPosition()
	g.now = time.Now()
	if g.mode == ModeTitle {
		g.updateNameField()
		g.updateCharacterSelect()
		g.updateRoomButtons()
		return nil
	}
	if g.isUndo() {
		g.Exit()
		g.mode = ModeTitle
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
		g.drawRoomButtons(Screen)
		return
	}
	g.drawRoomName(Screen)
	g.drawUndo(Screen)
	g.drawCharacters(Screen)
	g.drawSpeechBubble(Screen)
	g.drawMessageArea(Screen)
	g.drawMessageField(Screen)
}
