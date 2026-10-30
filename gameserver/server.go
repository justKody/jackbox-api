package engine

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/justKody/jackbox-api/types"
)

type GameServer struct {
	ctx      *actor.Context
	mu       sync.Mutex
	sessions map[int]*actor.PID
}

func newGameServer() actor.Producer {
	p.serverPID
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
	case types.PlayerAction:
		s.handlePlayerAction(msg)
	}
}

func (s *GameServer) handlePlayerAction(action types.PlayerAction) {

	switch action {

	case types.FindMatch:
		// Handle matchmaking


	case types.CancelMatch:

		// Cancel matchmaking

	case types.JoinRoom:

		// Join room

	case types.LeaveRoom:

		// Leave room

	default:
		// Unknown action
		fmt.Println("unkown player action registered in game server")

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

	// adding spawn child -> player session

	fmt.Println("new client trying connect")
	sid := rand.Intn(math.MaxInt)
	pid := s.ctx.SpawnChild(newPlayerSession(sid, conn, s.ctx.PID()), fmt.Sprintf("session_%s", sid))

	s.mu.Lock()
	s.sessions[sid] = pid
	s.mu.Unlock()
}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		fmt.Print("some error while starting the server:", err)
	}
	e.Spawn(newGameServer(), "server")
	select {}
}
