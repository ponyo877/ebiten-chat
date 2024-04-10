package drawable

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
	ScreenWidth        = 1400
	ScreenHeight       = 700
	MessageAreaPointX  = ScreenWidth * 0.7
	MessageFieldPointY = ScreenHeight * 0.95
	MessageAreaWidth   = ScreenWidth - MessageAreaPointX
	MessageFieldHeight = ScreenHeight - MessageFieldPointY
	NameFieldWidth     = ScreenWidth * 0.3
	NameFieldHeight    = ScreenHeight * 0.05
	NameFieldPointX    = ScreenWidth*0.5 - NameFieldWidth*0.5
	NameFieldPointY    = ScreenHeight*0.325 - NameFieldHeight*0.5
	RoomNamefontSize   = MessageFieldHeight * 3
	LargeFontSize      = MessageFieldHeight * 1.25
	MiddleFontSize     = MessageFieldHeight * 0.75
	SmallFontSize      = MessageFieldHeight * 0.5
	LogFontSize        = SmallFontSize * 0.8
)

var (
	CharacterImages  []*ebiten.Image
	RoomButtons      []*Button
	ArcadeFaceSource *text.GoTextFaceSource
	characterWidth   float64
	characterHeight  float64
)

var (
	RoomNameList = []string{
		"せり",
		"なずな",
		"ごぎょう",
		"はこべら",
		"ほとけのざ",
		"すずな",
		"すずしろ",
	}
)

func init() {
	CharacterImages = make([]*ebiten.Image, len(resources.Images))
	for i, img := range resources.Images {
		img, _, err := image.Decode(bytes.NewReader(img))
		if err != nil {
			log.Fatal(err)
		}
		CharacterImages[i] = ebiten.NewImageFromImage(img)
	}
	characterWidth = float64(CharacterImages[0].Bounds().Dx())
	characterHeight = float64(CharacterImages[0].Bounds().Dy())
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	ArcadeFaceSource = s
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}
