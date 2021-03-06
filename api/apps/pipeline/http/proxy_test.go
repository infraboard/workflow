package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/infraboard/mcube/logger/zap"
)

func TestWebsocket(t *testing.T) {
	// URL
	target := "ws://127.0.0.1:9948/workflow/api/v1/websocket/pipelines/c4s8iiea0brugin5hl30/watch"
	t.Logf("connnect to: %s", target)

	// Connect to the server
	h := http.Header{}
	h.Add("Authorization", "R9bRjrYpAtunMM9VDUPhCIgL")

	ws, _, err := websocket.DefaultDialer.Dial(target, h)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	fmt.Println(string(p))
}

func init() {
	zap.DevelopmentSetup()
}
