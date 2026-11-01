package engine

import (
	"fmt"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/justKody/jackbox-api/actorutil"
	"github.com/justKody/jackbox-api/game/rank"
	"github.com/justKody/jackbox-api/roomcode"
	"github.com/justKody/jackbox-api/types"
)

const (
	maxLobbyPlayers = 5               // the app's table is laid out for five
	matchmakingWait = 3 * time.Second // how long a player waits before being put in a lobby
)

// a player waiting out matchmakingWait before being placed in a lobby
type matchTicket struct {
	player  *actor.PID
	profile types.Profile
	timer   *time.Timer
}

// server side view of a lobby, used for matchmaking decisions
type lobbyInfo struct {
	id      string
	pid     *actor.PID
	players map[string]*actor.PID
	started bool // a game is running, nobody else gets placed here
}

// findMatch returns the open lobby with the most players, so lobbies fill up
// before new ones are created. Returns nil when every lobby is full or playing.
func (s *GameServer) findMatch() *lobbyInfo {
	var best *lobbyInfo
	for _, l := range s.lobbies {
		if l.started || len(l.players) >= maxLobbyPlayers {
			continue
		}
		if best == nil || len(l.players) > len(best.players) {
			best = l
		}
	}
	return best
}

// handleMatchmaking makes the player wait matchmakingWait, then places them in a lobby
func (s *GameServer) handleMatchmaking(player *actor.PID, profile types.Profile) {
	key := player.String()
	if lobbyId, ok := s.playerLobby[key]; ok {
		fmt.Println("player already in lobby", lobbyId)
		return
	}
	if _, ok := s.waiting[key]; ok {
		fmt.Println("player already waiting for a match")
		return
	}

	t := &matchTicket{player: player, profile: profile}
	t.timer = actorutil.SendAfter(s.ctx.Engine(), s.ctx.PID(), matchmakingWait, &types.ServerActionMessage{
		Action: types.PlacePlayer,
		Msg:    t,
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
	s.placeInLobby(t.player, t.profile)
}

func (s *GameServer) placeInLobby(player *actor.PID, profile types.Profile) {
	key := player.String()
	l := s.findMatch()
	if l == nil {
		l = s.createLobby()
	}

	l.players[key] = player
	s.playerLobby[key] = l.id
	s.ctx.Send(l.pid, &addPlayer{pid: player, profile: profile})
}

// handleBounced places a player again whose lobby started before they got in
func (s *GameServer) handleBounced(msg *playerBounced) {
	key := msg.pid.String()
	if l, ok := s.lobbies[s.playerLobby[key]]; ok {
		delete(l.players, key)
	}
	delete(s.playerLobby, key)
	s.placeInLobby(msg.pid, msg.profile)
}

// handleLobbyFinished settles a ranked game and frees its players to search again
func (s *GameServer) handleLobbyFinished(msg *lobbyFinished) {
	if msg.ranked {
		s.settleRanked(msg)
	}
	l, ok := s.lobbies[msg.lobbyId]
	if !ok {
		return
	}
	for key := range l.players {
		delete(s.playerLobby, key)
	}
	delete(s.lobbies, msg.lobbyId)
	s.ctx.Engine().Poison(l.pid)
}

// settleRanked moves everyone's RP by where they finished and tells each player their change
func (s *GameServer) settleRanked(msg *lobbyFinished) {
	seats := make(map[string]*seat, len(msg.seats)) // player id -> seat
	ratings := make(map[string]int, len(msg.seats)) // player id -> RP
	for _, st := range msg.seats {
		seats[st.player.ID] = st
		ratings[st.player.ID] = s.rating(st)
	}
	for _, r := range rank.Apply(msg.players, ratings) {
		st := seats[r.PlayerID]
		s.ratings[ratingKey(st)] = r.After
		s.ctx.Send(st.pid, &r)
	}
}

// rating is the player's RP, the first time a client shows up it brings the RP its app saved
func (s *GameServer) rating(st *seat) int {
	if rp, ok := s.ratings[ratingKey(st)]; ok {
		return rp
	}
	if st.profile.RP > 0 {
		return st.profile.RP
	}
	return rank.StartingRP
}

// ratings follow the app's client id, a client without one only keeps RP for this connection
func ratingKey(st *seat) string {
	if st.profile.ClientID != "" {
		return st.profile.ClientID
	}
	return st.pid.String()
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
	id := roomcode.New(func(code string) bool {
		_, taken := s.lobbies[code]
		return taken
	})
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
