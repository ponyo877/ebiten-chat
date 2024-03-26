//go:build js && wasm
// +build js,wasm

package main

import (
	_ "embed"
	"flag"
	_ "image/png"
	"syscall/js"

	"github.com/ponyo877/folks-ui/game"
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
	g := game.NewGame(wsScheme, wsHost)

	// ブラウザを閉じた時の処理
	onBeforeunload = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.Exit()
		return nil
	})
	window.Call("addEventListener", "beforeunload", onBeforeunload)

	if err := g.Start(); err != nil {
		panic(err)
	}
	if err := g.Exit(); err != nil {
		panic(err)
	}
}
