LOCALLDFLAGS := -X 'main.wsScheme=ws' \
                -X 'main.wsHost=localhost:8000'
PRDLDFLAGS := -X 'main.wsScheme=wss' \
              -X 'main.wsHost=folks-chat.com'

localwasm: main.go
	GOOS=js GOARCH=wasm go build -ldflags "$(LOCALLDFLAGS)" -o ebiten_chat.wasm $<

prdwasm: main.go
	GOOS=js GOARCH=wasm go build -ldflags "$(PRDLDFLAGS)" -o ebiten_chat.wasm $<
