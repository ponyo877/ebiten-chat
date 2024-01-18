package folks

import (
	"bytes"
	"flag"
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	resources "github.com/ponyo877/folks-ui/static"
)

var FlagCRT = flag.Bool("crt", false, "enable the CRT effect")

const (
	ScreenWidth      = 640
	ScreenHeight     = 480
	tileSize         = 32
	titleFontSize    = fontSize * 1.5
	fontSize         = 24
	smallFontSize    = fontSize / 2
	pipeWidth        = tileSize * 2
	pipeStartOffsetX = 8
	pipeIntervalX    = 8
	pipeGapY         = 5
)

var (
	gopherLeftImage  *ebiten.Image
	gopherRightImage *ebiten.Image
	arcadeFaceSource *text.GoTextFaceSource
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.Gopher_Left_png))
	if err != nil {
		log.Fatal(err)
	}
	gopherLeftImage = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(resources.Gopher_Right_png))
	if err != nil {
		log.Fatal(err)
	}
	gopherRightImage = ebiten.NewImageFromImage(img)
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = s
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}
