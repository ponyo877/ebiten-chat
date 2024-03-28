package game

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	d "github.com/ponyo877/folks-ui/drawable"
)

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
				clr = color.RGBA{0, 0, 255, 255}
			}
			if i == g.clickedX && j == g.clickedY {
				clr = color.RGBA{255, 0, 0, 255}
			}
			if clr != color.White {
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
