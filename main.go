package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	_ "image/png"
	"io"
	"net/http"

	"github.com/ponyo877/folks-ui/folks"

	"github.com/hajimehoshi/ebiten/v2"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	response := Response{}
	flag.Parse()
	ebiten.SetWindowSize(folks.ScreenWidth, folks.ScreenHeight)
	resp, _ := http.Get("https://folks-chat.com/v1/")
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return
	}
	// ebiten.SetWindowTitle("Flappy Gopher (Ebitengine Demo)")
	ebiten.SetWindowTitle(response.Message)
	if err := ebiten.RunGame(folks.NewGame(*folks.FlagCRT)); err != nil {
		panic(err)
	}
}
