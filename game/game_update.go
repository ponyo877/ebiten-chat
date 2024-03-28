package game

import (
	"math/rand"
	"slices"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	d "github.com/ponyo877/folks-ui/drawable"
	"github.com/ponyo877/folks-ui/entity"
)

func (g *Game) updateNameField() {
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
}

func (g *Game) updateCharacterSelect() {
	if g.x-startSelectX > 0 && g.y-startSelectY > 0 {
		g.bluredX, g.bluredY = (g.x-startSelectX)/cellSize, (g.y-startSelectY)/cellSize
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.clickedX, g.clickedY = g.bluredX, g.bluredY
		imgid := g.clickedX + g.clickedY*NumOfImagesPerRow
		if 0 <= imgid && imgid < len(d.CharacterImage) {
			g.imgid = g.clickedX + g.clickedY*NumOfImagesPerRow
		}
	}
}

func (g *Game) updateEnterButton() {
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

func (g *Game) updateCharacterDir() {
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
}

func (g *Game) updateChatMsg() {
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
}

func (g *Game) updateMsgField() {
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
}

func (g *Game) updateMove() {
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
