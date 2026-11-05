package gameserver

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/justKody/jackbox-api/internal/types"
	"github.com/justKody/jackbox-api/internal/wsmsg"
)

type PlayerSession struct {
	conn      *websocket.Conn
	ctx       *actor.Context
	serverPID *actor.PID
}

func newPlayerSession(conn *websocket.Conn, serverPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &PlayerSession{
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
	case *types.GameUpdate:
		p.writeJSON("game", msg)
	case *types.DialUpdate:
		p.writeJSON("dial", msg)
	case *types.ActionError:
		p.writeJSON("error", msg)
	default:
		fmt.Println("msg received for player session", msg)
	}
}

func (p *PlayerSession) writeJSON(msgType string, data any) {
	if err := wsmsg.Write(p.conn, msgType, data); err != nil {
		fmt.Println("write error in player session", err)
	}
}

func (p *PlayerSession) readLoop() {
	for {
		var msg types.WSMessage
		if err := p.conn.ReadJSON(&msg); err != nil {
			fmt.Println("read error in read loop", err)
			// a close frame or a dropped connection both mean the client is gone,
			// free up its lobby slot and stop this session
			p.sendAction(types.CancelMatch, nil)
			p.conn.Close()
			p.ctx.Engine().Poison(p.ctx.PID())
			return
		}

		p.handleMessage(msg)
	}
}

func (p *PlayerSession) handleMessage(msg types.WSMessage) {
	switch msg.Type {
	case "findMatch":
		// ask server to put urself in a finding match logic
		var profile types.Profile
		if !p.decode(msg, &profile) {
			return
		}
		p.sendAction(types.FindMatch, profile)
	case "cancelMatch":
		p.sendAction(types.CancelMatch, nil)
	case "lockGuess":
		var m types.LockGuessMsg
		if p.decode(msg, &m) {
			p.sendAction(types.LockGuess, m)
		}
	case "vote":
		var m types.VoteMsg
		if p.decode(msg, &m) {
			p.sendAction(types.Vote, m)
		}
	case "dial":
		var m types.DialMsg
		if p.decode(msg, &m) {
			p.sendAction(types.Dial, m)
		}
	default:
		fmt.Println("msg rec in player session:", msg)
	}
}

// decode reads a message's data, logging and dropping the message when it's malformed
func (p *PlayerSession) decode(msg types.WSMessage, v any) bool {
	if err := wsmsg.Decode(msg, v); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// ctx.Send sets this session as the sender, which the server uses to identify the player
func (p *PlayerSession) sendAction(action types.PlayerAction, msg any) {
	p.ctx.Send(p.serverPID, &types.PlayerActionMessage{
		Action: action,
		Msg:    msg,
	})
}
