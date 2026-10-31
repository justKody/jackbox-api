package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/justKody/jackbox-api/types"
)

type PlayerSession struct {
	sessionId int
	username  string
	// clientId -> only use for redis re connection thing
	inGame    bool // do we need this?
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
		// ctx must be set before readLoop sends anything with it
		p.ctx = c
		go p.readLoop()

	// writes only happen here, so the websocket never sees concurrent writers
	case *types.LobbyUpdate:
		p.writeJSON("lobbyUpdate", msg)
	case *types.MatchCancelled:
		p.writeJSON("matchCancelled", msg)
	default:
		fmt.Println("msg received for player session", msg)
	}
}

// write this out somewhere else ( maybe util)
func (p *PlayerSession) writeJSON(msgType string, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("marshal error in player session", err)
		return
	}
	if err := p.conn.WriteJSON(types.WSMessage{Type: msgType, Data: b}); err != nil {
		fmt.Println("write error in player session", err)
	}
}

func (p *PlayerSession) readLoop() {
	for {
		var msg types.WSMessage
		if err := p.conn.ReadJSON(&msg); err != nil {
			fmt.Println("read error in read loop", err)
			if errors.Is(err, io.EOF) {
				// client is gone, free up its lobby slot
				p.sendAction(types.CancelMatch)
			}
			return
		}

		p.handleMessage(msg)
	}
}

func (p *PlayerSession) handleMessage(msg types.WSMessage) {
	switch msg.Type {
	case "findMatch":
		// ask server to put urself in a finding match logic
		p.sendAction(types.FindMatch)
	case "cancelMatch":
		p.sendAction(types.CancelMatch)
	default:
		fmt.Println("msg rec in player session:", msg)
	}
}

// ctx.Send sets this session as the sender, which the server uses to identify the player
func (p *PlayerSession) sendAction(action types.PlayerAction) {
	p.ctx.Send(p.serverPID, &types.PlayerActionMessage{
		Action: action,
		Msg:    nil,
	})
}
