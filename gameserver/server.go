package engine

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
)

type GameServer struct {
	ctx      *actor.Context
	sessions map[int]*actor.PID
}

func newGameServer() actor.Producer {
	return func() actor.Receiver {
		return &GameServer{
			sessions: make(map[int]*actor.PID),
		}
	}
}

func (s *GameServer) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		// start the server handling websocket
		s.startHttp()
		s.ctx = c
		_ = msg
	}
}

func (s *GameServer) startHttp() {
	fmt.Println("starting HTTP server on port 4000")
	// replace it with mux router
	go func() {
		http.HandleFunc("/ws", s.handleWs)
		http.ListenAndServe(":4000", nil)
	}()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		// allowing no origin header and any localhost dev origin
		return origin == "" || strings.HasPrefix(origin, "http://localhost")
	},
}

func (s *GameServer) handleWs(w http.ResponseWriter, r *http.Request) {
	// add authentication check here

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return // Upgrade already wrote the HTTP error
	}

}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		fmt.Print("some error while starting the server:", err)
	}
	e.Spawn(newGameServer(), "server")
	select {}
}
