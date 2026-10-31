package engine

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/justKody/jackbox-api/types"
)

const (
	maxLobbyPlayers = 8
	lobbyCodeLength = 4
	lobbyCodeChars  = "ABCDEFGHJKLMNPQRSTUVWXYZ" // no I/O so codes are easy to read out loud
	matchmakingWait = 3 * time.Second            // how long a player waits before being put in a lobby
)

// a player waiting out matchmakingWait before being placed in a lobby
type matchTicket struct {
	player *actor.PID
	timer  *time.Timer
}

// server side view of a lobby, used for matchmaking decisions
type lobbyInfo struct {
	id      string
	pid     *actor.PID
	players map[string]*actor.PID
}

// findMatch returns the open lobby with the most players, so lobbies fill up
// before new ones are created. Returns nil when every lobby is full.
func (s *GameServer) findMatch() *lobbyInfo {
	var best *lobbyInfo
	for _, l := range s.lobbies {
		if len(l.players) >= maxLobbyPlayers {
			continue
		}
		if best == nil || len(l.players) > len(best.players) {
			best = l
		}
	}
	return best
}

// handleMatchmaking makes the player wait matchmakingWait, then places them in a lobby
func (s *GameServer) handleMatchmaking(player *actor.PID) {
	key := player.String()
	if lobbyId, ok := s.playerLobby[key]; ok {
		fmt.Println("player already in lobby", lobbyId)
		return
	}
	if _, ok := s.waiting[key]; ok {
		fmt.Println("player already waiting for a match")
		return
	}

	t := &matchTicket{player: player}
	// the timer runs on its own goroutine, so it reports back via a message
	// instead of touching server state directly
	engine, self := s.ctx.Engine(), s.ctx.PID()
	t.timer = time.AfterFunc(matchmakingWait, func() {
		engine.Send(self, &types.ServerActionMessage{
			Action: types.PlacePlayer,
			Msg:    t,
		})
	})
	s.waiting[key] = t
}

// handlePlacePlayer puts the player in an open lobby, creating one if needed
func (s *GameServer) handlePlacePlayer(t *matchTicket) {
	key := t.player.String()
	// a ticket that was cancelled (or replaced by a newer one) is stale
	if s.waiting[key] != t {
		return
	}
	delete(s.waiting, key)
	player := t.player

	l := s.findMatch()
	if l == nil {
		l = s.createLobby()
	}

	l.players[key] = player
	s.playerLobby[key] = l.id
	s.ctx.Send(l.pid, &addPlayer{pid: player})
}

// handleCancelMatch stops the player's wait, or takes them out of their lobby
// and closes the lobby once empty
func (s *GameServer) handleCancelMatch(player *actor.PID) {
	key := player.String()
	if t, ok := s.waiting[key]; ok {
		t.timer.Stop()
		delete(s.waiting, key)
		return
	}

	lobbyId, ok := s.playerLobby[key]
	if !ok {
		return
	}
	l := s.lobbies[lobbyId]

	delete(s.playerLobby, key)
	delete(l.players, key)
	s.ctx.Send(l.pid, &removePlayer{pid: player})

	if len(l.players) == 0 {
		delete(s.lobbies, lobbyId)
		s.ctx.Engine().Poison(l.pid)
	}
}

// spawns the child, used if there exist no lobby
func (s *GameServer) createLobby() *lobbyInfo {
	id := s.newLobbyCode()
	pid := s.ctx.SpawnChild(newLobby(s.ctx.PID(), id), fmt.Sprintf("lobby_%s", id))
	l := &lobbyInfo{
		id:      id,
		pid:     pid,
		players: make(map[string]*actor.PID),
	}
	s.lobbies[id] = l
	fmt.Println("created lobby", id)
	return l
}

// jackbox style room code, e.g. "QXRT"
func (s *GameServer) newLobbyCode() string {
	for {
		b := make([]byte, lobbyCodeLength)
		for i := range b {
			b[i] = lobbyCodeChars[rand.Intn(len(lobbyCodeChars))]
		}
		if _, taken := s.lobbies[string(b)]; !taken {
			return string(b)
		}
	}
}
