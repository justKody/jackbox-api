package engine

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/justKody/jackbox-api/types"
)

type Lobby struct {
	lobbyId   string
	players   map[string]*actor.PID // player session pid -> pid
	ctx       *actor.Context
	serverPID *actor.PID
}

func newLobby(serverPID *actor.PID, lobbyId string) actor.Producer {
	return func() actor.Receiver {
		return &Lobby{
			lobbyId:   lobbyId,
			serverPID: serverPID,
			players:   make(map[string]*actor.PID),
		}
	}
}

// sent by the game server to a lobby when matchmaking places a player in it
type addPlayer struct {
	pid *actor.PID
}

// sent by the game server to a lobby when a player cancels or disconnects
type removePlayer struct {
	pid *actor.PID
}

func (l *Lobby) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		l.ctx = c
	case *addPlayer:
		l.players[msg.pid.String()] = msg.pid
		l.broadcastUpdate()
	case *removePlayer:
		delete(l.players, msg.pid.String())
		c.Send(msg.pid, &types.MatchCancelled{LobbyID: l.lobbyId})
		l.broadcastUpdate()
	case actor.Stopped:
		fmt.Println("lobby closed:", l.lobbyId)
	}
}

// tell every player in the lobby the current roster size
func (l *Lobby) broadcastUpdate() {
	update := &types.LobbyUpdate{
		LobbyID:    l.lobbyId,
		Players:    len(l.players),
		MaxPlayers: maxLobbyPlayers,
	}
	for _, pid := range l.players {
		l.ctx.Send(pid, update)
	}
}
