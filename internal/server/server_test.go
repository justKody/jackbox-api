package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/justKody/jackbox-api/internal/ws"
	"github.com/justKody/jackbox-api/internal/ws/types"
)

// the real routes in front of a real game server
func startServer(t *testing.T) (*ws.Service, *httptest.Server) {
	t.Helper()
	games, err := ws.Start()
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(NewServer(":0", "../../s3asset", games.HandleWs).Router)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		games.Stop(ctx)
		ts.Close()
	})
	return games, ts
}

func dial(t *testing.T, ts *httptest.Server) *websocket.Conn {
	t.Helper()
	conn, resp, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(ts.URL, "http")+"/ws", nil)
	if err != nil {
		t.Fatalf("dial /ws: %v", err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("dial /ws: status %d", resp.StatusCode)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestRoutes(t *testing.T) {
	_, ts := startServer(t)

	for path, want := range map[string]int{
		"/v1/health":           http.StatusOK,
		"/assets/v2/owid1.mp3": http.StatusOK,
		// only the voice-overs are served, not the rest of s3asset
		"/assets/voicelines.md": http.StatusNotFound,
		"/nope":                 http.StatusNotFound,
	} {
		resp, err := http.Get(ts.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != want {
			t.Errorf("GET %s = %d, want %d", path, resp.StatusCode, want)
		}
	}
}

func TestFindMatchPutsClientInLobby(t *testing.T) {
	_, ts := startServer(t)
	conn := dial(t, ts)

	data, _ := json.Marshal(types.Profile{Name: "Ada"})
	if err := conn.WriteJSON(types.WSMessage{Type: "findMatch", Data: data}); err != nil {
		t.Fatal(err)
	}

	// matchmaking makes a player wait a few seconds
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	var msg types.WSMessage
	if err := conn.ReadJSON(&msg); err != nil {
		t.Fatal(err)
	}
	if msg.Type != "lobbyUpdate" {
		t.Fatalf("got message %q, want lobbyUpdate", msg.Type)
	}
	var update types.LobbyUpdate
	if err := json.Unmarshal(msg.Data, &update); err != nil {
		t.Fatal(err)
	}
	if len(update.Players) != 1 || update.Players[0].Name != "Ada" {
		t.Errorf("players = %+v, want just Ada", update.Players)
	}
}

// stopping the game server has to close websockets, http.Server.Shutdown won't
func TestStopClosesWebsockets(t *testing.T) {
	games, ts := startServer(t)
	conn := dial(t, ts)

	// a session only exists once the actor has handled the connection, a read
	// deadline keeps a missed close from hanging the test
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	games.Stop(ctx)
	if ctx.Err() != nil {
		t.Fatal("Stop timed out")
	}

	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Fatal("read succeeded after Stop")
	}
	if ne, ok := err.(interface{ Timeout() bool }); ok && ne.Timeout() {
		t.Fatalf("connection still open after Stop: %v", err)
	}
}
