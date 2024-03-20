//go:build js && wasm
// +build js,wasm

package main

import (
	_ "embed"
	"flag"
	_ "image/png"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ponyo877/folks-ui/folks"
)

var (
	wsScheme string
	wsHost   string
)

func main() {
	var onBeforeunload js.Func
	defer onBeforeunload.Release()
	window := js.Global()

	flag.Parse()
	ebiten.SetWindowSize(folks.ScreenWidth, folks.ScreenHeight)
	ebiten.SetWindowTitle("Ebiten Chat")
	game := folks.NewGame(wsScheme, wsHost)

	// ブラウザを閉じた時の処理
	onBeforeunload = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		game.Exit()
		return nil
	})
	window.Call("addEventListener", "beforeunload", onBeforeunload)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
	if err := game.Exit(); err != nil {
		panic(err)
	}
}
