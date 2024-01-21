package main

import (
	_ "embed"
	"flag"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ponyo877/folks-ui/folks"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	response := Response{}
	flag.Parse()
	ebiten.SetWindowSize(folks.ScreenWidth, folks.ScreenHeight)
	ebiten.SetWindowTitle("Ebiten Chat")
	ebiten.SetWindowTitle(response.Message)
	game := folks.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
	if err := game.Exit(); err != nil {
		panic(err)
	}
}
