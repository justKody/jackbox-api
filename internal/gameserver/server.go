package gameserver

import (
	"fmt"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/justKody/jackbox-api/internal/types"
	"github.com/justKody/jackbox-api/internal/wsmsg"
)

type GameServer struct {
	ctx *actor.Context

	// only touched inside Receive, so no lock needed
	lobbies     map[string]*lobbyInfo   // lobby id -> lobby
	playerLobby map[string]string       // player session pid -> lobby id
	waiting     map[string]*matchTicket // player session pid -> ticket, while waiting for a lobby
}

// sent by HandleWs when a client's websocket is ready for a player session
type clientConnected struct {
	conn *websocket.Conn
}

func newGameServer() actor.Producer {
	return func() actor.Receiver {
		return &GameServer{
			lobbies:     make(map[string]*lobbyInfo),
			playerLobby: make(map[string]string),
			waiting:     make(map[string]*matchTicket),
		}
	}
}

func (s *GameServer) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		s.ctx = c
	case *clientConnected:
		fmt.Println("new client connected")
		c.SpawnChild(newPlayerSession(msg.conn, c.PID()), "session")
	case *types.PlayerActionMessage:
		s.handlePlayerAction(c.Sender(), msg)
	case *types.ServerActionMessage:
		s.handleServerAction(msg)
	case *lobbyStarted:
		if l, ok := s.lobbies[msg.lobbyId]; ok {
			l.started = true
		}
	case *playerBounced:
		s.handleBounced(msg)
	case *lobbyFinished:
		s.handleLobbyFinished(msg)
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

func (s *GameServer) handlePlayerAction(player *actor.PID, msg *types.PlayerActionMessage) {

	switch msg.Action {

	case types.FindMatch:
		profile, _ := msg.Msg.(types.Profile)
		s.handleMatchmaking(player, profile)

	case types.CancelMatch:
		s.handleCancelMatch(player)

	case types.JoinRoom:
		// Join room (not needed right now)
	case types.LeaveRoom:
		// Leave room (now needed right now)

	case types.LockGuess, types.Vote, types.Dial:
		// game moves go straight to the player's lobby
		if l, ok := s.lobbies[s.playerLobby[player.String()]]; ok {
			s.ctx.Send(l.pid, &playerMove{pid: player, action: msg.Action, msg: msg.Msg})
		}

	default:
		// Unknown action
		fmt.Println("unkown player action registered in game server")

	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return wsmsg.AllowOrigin(r.Header.Get("Origin"), r.Host)
	},
}

// Service is the running game server as the rest of the program sees it
type Service struct {
	engine    *actor.Engine
	serverPID *actor.PID
}

// Start runs the actor engine and the game server actor
func Start() (*Service, error) {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		return nil, fmt.Errorf("could not start the actor engine: %w", err)
	}
	return &Service{engine: e, serverPID: e.Spawn(newGameServer(), "server")}, nil
}

// HandleWs upgrades a client to a websocket and hands it to the game server actor.
// It runs on an http goroutine, so it only talks to the actor with a message
func (g *Service) HandleWs(w http.ResponseWriter, r *http.Request) {
	// add authentication check here

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return // Upgrade already wrote the HTTP error
	}
	g.engine.Send(g.serverPID, &clientConnected{conn: conn})
}
