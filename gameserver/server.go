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

	// only touched inside Receive, so no lock needed
	lobbies     map[string]*lobbyInfo   // lobby id -> lobby
	playerLobby map[string]string       // player session pid -> lobby id
	waiting     map[string]*matchTicket // player session pid -> ticket, while waiting for a lobby
}

func newGameServer() actor.Producer {
	return func() actor.Receiver {
		return &GameServer{
			sessions:    make(map[int]*actor.PID),
			lobbies:     make(map[string]*lobbyInfo),
			playerLobby: make(map[string]string),
			waiting:     make(map[string]*matchTicket),
		}
	}
}

func (s *GameServer) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		// ctx must be set before the http server can spawn sessions with it
		s.ctx = c
		s.startHttp()
	case *types.PlayerActionMessage:
		s.handlePlayerAction(c.Sender(), msg.Action)
	case *types.ServerActionMessage:
		s.handleServerAction(msg)
	}
}

func (s *GameServer) handleServerAction(msg *types.ServerActionMessage) {
	switch msg.Action {
	case types.PlacePlayer:
		s.handlePlacePlayer(msg.Msg.(*matchTicket))
	default:
		fmt.Println("unknown server action registered in game server")
	}
}

func (s *GameServer) handlePlayerAction(player *actor.PID, action types.PlayerAction) {

	switch action {

	case types.FindMatch:
		s.handleMatchmaking(player)

	case types.CancelMatch:
		s.handleCancelMatch(player)

	case types.JoinRoom:

		// Join room (not needed right now)

	case types.LeaveRoom:

		// Leave room (now needed right now)

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
	pid := s.ctx.SpawnChild(newPlayerSession(sid, conn, s.ctx.PID()), fmt.Sprintf("session_%d", sid))

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
