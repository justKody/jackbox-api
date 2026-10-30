package engine

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/justKody/jackbox-api/types"
)

type PlayerSession struct {
	sessionId int
	username  string
	// clientId -> only use for redis re connection thing
	inGame    bool
	conn      *websocket.Conn
	ctx       *actor.Context
	serverPID *actor.PID
}

func newPlayerSession(sessionId int, conn *websocket.Conn, serverPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &PlayerSession{
			sessionId: sessionId,
			inGame:    false,
			conn:      conn,
			serverPID: serverPID,
		}
	}
}

func (p *PlayerSession) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		go p.readLoop()

	// case Message: .. a case of state represented as method calling case
	// 	// sending json to this player
	default:
		fmt.Println("msg received for player session", msg)
	}
}

func (p *PlayerSession) readLoop() {
	var msg types.WSMessage
	for {
		if err := p.conn.ReadJSON(&msg); err != nil {
			fmt.Println("read error in read loop", err)
		}

		go p.handleMessage(msg)
	}
}

func (p *PlayerSession) handleMessage(msg types.WSMessage) {
	switch msg.Type {
	case "findMatch":
		// logic goes here
	default:
		fmt.Println("msg rec in player session:", msg)
	}
}
