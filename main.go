//go:build js && wasm
// +build js,wasm

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
	// var onBeforeunload js.Func
	// defer onBeforeunload.Release()

	response := Response{}
	flag.Parse()
	ebiten.SetWindowSize(folks.ScreenWidth, folks.ScreenHeight)
	ebiten.SetWindowTitle("Ebiten Chat")
	ebiten.SetWindowTitle(response.Message)
	game := folks.NewGame()

	// ブラウザを閉じた時の処理
	// onBeforeunload = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	game.Exit()
	// 	return nil
	// })
	// js.Global().Call("addEventListener", "beforeunload", onBeforeunload)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
	if err := game.Exit(); err != nil {
		panic(err)
	}
}
