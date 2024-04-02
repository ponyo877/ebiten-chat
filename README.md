
# Ebiten Chat
Ebiten Chat is a simple avatar chat created with Ebitengine.
DEMO: https://folks-chat.com/

# Features
- Written entirely in Go language.
- Real-time chatting through a WebSocket server.
- You can freely move your chosen avatar within the browser.

# Usage
You need to implement and connect to a WebSocket server.
The source code for the WebSocket server is not published, but it can be easily written in Go using gorilla/websocket.

For the client, build the WebAssembly and upload it to your web server.
```bash
GOOS=js GOARCH=wasm go build -ldflags -X 'main.wsScheme=ws -X main.wsHost=localhost:8000' -o ebiten_chat.wasm
```