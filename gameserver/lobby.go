package engine

import (
	"github.com/anthdm/hollywood/actor"
)

type Lobby struct {
	lobbyId   string
	players   map[*PlayerSession]bool
	ctx       *actor.Context
	serverPID *actor.PID
}

func newLobby(serverPID *actor.PID, lobbyId string) actor.Producer {
	return func() actor.Receiver {
		return &Lobby{
			lobbyId:   lobbyId,
			serverPID: serverPID,
			players:   make(map[*PlayerSession]bool),
		}
	}
}

func (l *Lobby) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		l.ctx = c
		_ = msg
	}
}
