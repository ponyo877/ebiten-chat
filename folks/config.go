package folks

import (
	"bytes"
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	resources "github.com/ponyo877/folks-ui/static"
)

const (
	ScreenWidth       = 1400
	ScreenHeight      = 700
	MessageAreaPointX = ScreenWidth * 0.7
	MessageAreaWidth  = ScreenWidth - MessageAreaPointX
	TextFieldPointY   = ScreenHeight * 0.95
	TextFieldHeight   = ScreenHeight - TextFieldPointY
	fontSize          = TextFieldHeight * 0.75
	smallFontSize     = TextFieldHeight * 0.5
	logFontSize       = smallFontSize * 0.8
)

var (
	characterImage   []*ebiten.Image
	arcadeFaceSource *text.GoTextFaceSource
	characterWidth   float64
	characterHeight  float64
)

func init() {
	characterImage = make([]*ebiten.Image, len(resources.Images))
	for i, img := range resources.Images {
		img, _, err := image.Decode(bytes.NewReader(img))
		if err != nil {
			log.Fatal(err)
		}
		characterImage[i] = ebiten.NewImageFromImage(img)
	}
	characterWidth = float64(characterImage[0].Bounds().Dx())
	characterHeight = float64(characterImage[0].Bounds().Dy())
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = s
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}
