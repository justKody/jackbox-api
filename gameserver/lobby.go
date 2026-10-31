package engine

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/justKody/jackbox-api/game"
	gametype "github.com/justKody/jackbox-api/game/types"
	"github.com/justKody/jackbox-api/types"
)

const (
	minLobbyPlayers = 2
	lobbyStartWait  = 10 * time.Second // countdown once enough players are in, skipped when the lobby fills up
	revealWait      = 15 * time.Second // the reveal moves on when everyone pressed next, or after this
	allVotedWait    = 900 * time.Millisecond
	// clients lock their own guess when their dial runs out, the server only steps in if that never arrives
	timerGrace = 2 * time.Second
)

// seat colors, in join order
var seatColors = []string{"#FFD23F", "#4C8DFF", "#FF8A3D", "#B58CFF", "#5ED3A8"}

type seat struct {
	pid     *actor.PID
	profile types.Profile
	player  gametype.Player
}

type Lobby struct {
	lobbyId   string
	seats     []*seat          // join order
	byPID     map[string]*seat // player session pid -> seat
	joined    int              // seats handed out so far, for seat ids
	ctx       *actor.Context
	serverPID *actor.PID

	startsAt   time.Time
	startTimer *time.Timer

	game      *game.Game
	abandoned bool            // ended early because players left
	phaseKey  string          // round/phase/guesser the current deadline belongs to
	tick      int             // bumped on every new deadline, so stale timers are ignored
	timer     *time.Timer     // the current phase's deadline
	revealing bool            // everyone voted, reveal is scheduled
	ready     map[string]bool // player id -> pressed next on the reveal
}

func newLobby(serverPID *actor.PID, lobbyId string) actor.Producer {
	return func() actor.Receiver {
		return &Lobby{
			lobbyId:   lobbyId,
			serverPID: serverPID,
			byPID:     make(map[string]*seat),
		}
	}
}

// sent by the game server to a lobby when matchmaking places a player in it
type addPlayer struct {
	pid     *actor.PID
	profile types.Profile
}

// sent by the game server to a lobby when a player cancels or disconnects
type removePlayer struct {
	pid *actor.PID
}

// sent by the game server to a lobby with a player's game action
type playerMove struct {
	pid    *actor.PID
	action types.PlayerAction
	msg    any
}

// sent by the lobby to itself from timers
type startGame struct{}
type phaseTimeout struct{ tick int }

// sent by a lobby to the game server
type lobbyStarted struct{ lobbyId string }

// a player placed in a lobby that already started, the server places them again
type playerBounced struct {
	pid     *actor.PID
	profile types.Profile
}

type lobbyFinished struct {
	lobbyId string
	// false when the game ended early, nobody's RP moves then
	ranked bool
	// still at the table with final scores, seats is who they are
	players []gametype.Player
	seats   []*seat
}

func (l *Lobby) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		l.ctx = c
	case *addPlayer:
		l.handleAdd(msg)
	case *removePlayer:
		l.handleRemove(msg.pid)
		c.Send(msg.pid, &types.MatchCancelled{LobbyID: l.lobbyId})
	case *startGame:
		l.start()
	case *playerMove:
		l.handleMove(msg)
	case *phaseTimeout:
		l.handleTimeout(msg.tick)
	case actor.Stopped:
		l.stopTimers()
		fmt.Println("lobby closed:", l.lobbyId)
	}
}

func (l *Lobby) handleAdd(msg *addPlayer) {
	if l.game != nil {
		l.ctx.Send(l.serverPID, &playerBounced{pid: msg.pid, profile: msg.profile})
		return
	}
	l.joined++
	name := msg.profile.Name
	if name == "" {
		name = fmt.Sprintf("Player %d", l.joined)
	}
	s := &seat{
		pid:     msg.pid,
		profile: msg.profile,
		player: gametype.Player{
			ID:       fmt.Sprintf("p%d", l.joined),
			Name:     name,
			Initials: initials(name),
			Color:    seatColors[(l.joined-1)%len(seatColors)],
			Animal:   msg.profile.Animal,
			Skin:     msg.profile.Skin,
		},
	}
	l.seats = append(l.seats, s)
	l.byPID[msg.pid.String()] = s

	switch {
	case len(l.seats) >= maxLobbyPlayers:
		l.start()
		return
	case len(l.seats) >= minLobbyPlayers && l.startTimer == nil:
		l.startsAt = time.Now().Add(lobbyStartWait)
		l.startTimer = l.after(lobbyStartWait, &startGame{})
	}
	l.broadcastLobby()
}

func (l *Lobby) handleRemove(pid *actor.PID) {
	s, ok := l.byPID[pid.String()]
	if !ok {
		return
	}
	delete(l.byPID, pid.String())
	for i, o := range l.seats {
		if o == s {
			l.seats = append(l.seats[:i], l.seats[i+1:]...)
			break
		}
	}

	if l.game == nil {
		if len(l.seats) < minLobbyPlayers && l.startTimer != nil {
			l.startTimer.Stop()
			l.startTimer = nil
		}
		l.broadcastLobby()
		return
	}
	if l.game.Phase == gametype.PhaseGameOver {
		return
	}
	if !l.game.RemovePlayer(s.player.ID) {
		l.game.Phase = gametype.PhaseGameOver
		l.abandoned = true
	}
	delete(l.ready, s.player.ID)
	if l.game.Phase == gametype.PhaseReveal && l.allReady() {
		l.game.Next()
	}
	l.sync()
}

func (l *Lobby) start() {
	if l.game != nil || len(l.seats) < minLobbyPlayers {
		return
	}
	if l.startTimer != nil {
		l.startTimer.Stop()
		l.startTimer = nil
	}
	players := make([]gametype.Player, len(l.seats))
	for i, s := range l.seats {
		players[i] = s.player
	}
	g, err := game.NewGame(players, nil)
	if err != nil {
		fmt.Println("could not start game in lobby", l.lobbyId, err)
		return
	}
	l.game = g
	l.ctx.Send(l.serverPID, &lobbyStarted{lobbyId: l.lobbyId})
	fmt.Println("game started in lobby", l.lobbyId)
	l.sync()
}

func (l *Lobby) handleMove(m *playerMove) {
	s, ok := l.byPID[m.pid.String()]
	if !ok || l.game == nil {
		return
	}
	id := s.player.ID
	var err error
	switch m.action {
	case types.LockGuess:
		err = l.game.LockGuess(id, m.msg.(types.LockGuessMsg).Guess)
	case types.Vote:
		err = l.game.Vote(id, m.msg.(types.VoteMsg).Call)
	case types.Dial:
		// the dial streams in while the guesser drags it, so it skips the full game update.
		// one arriving just after the lock is expected, it's dropped without an error
		if l.game.MoveDial(id, m.msg.(types.DialMsg).Value) == nil {
			l.broadcastDial(s)
		}
		return
	case types.Next:
		if l.game.Phase != gametype.PhaseReveal {
			err = game.ErrWrongPhase
			break
		}
		l.ready[id] = true
		if l.allReady() {
			err = l.game.Next()
		}
	}
	if err != nil {
		l.ctx.Send(m.pid, &types.ActionError{Message: err.Error()})
		return
	}
	l.sync()
}

// a phase ran out of time, do what the clients didn't
func (l *Lobby) handleTimeout(tick int) {
	if tick != l.tick || l.game == nil {
		return
	}
	switch l.game.Phase {
	case gametype.PhaseGuess:
		l.game.GuessTimeout()
	case gametype.PhaseVote:
		l.game.Reveal()
	case gametype.PhaseReveal:
		l.game.Next()
	}
	l.sync()
}

// sync arms the deadline for the phase the game is now in and sends everyone the new state
func (l *Lobby) sync() {
	g := l.game
	key := fmt.Sprintf("%d/%s/%s", g.Round, g.Phase, g.Guesser().ID)
	if key != l.phaseKey {
		l.phaseKey = key
		l.revealing = false
		switch g.Phase {
		case gametype.PhaseGuess:
			l.deadline(game.GuessSeconds + timerGrace)
		case gametype.PhaseVote:
			l.deadline(game.VoteSeconds + timerGrace)
		case gametype.PhaseReveal:
			l.ready = make(map[string]bool)
			l.deadline(revealWait)
		case gametype.PhaseGameOver:
			l.stopTimers()
		}
	}
	// reveal shortly after the last call comes in
	if g.Phase == gametype.PhaseVote && !l.revealing && g.AllVoted() {
		l.revealing = true
		l.deadline(allVotedWait)
	}

	l.broadcastGame()
	if g.Phase == gametype.PhaseGameOver {
		seats := append([]*seat(nil), l.seats...)
		l.ctx.Send(l.serverPID, &lobbyFinished{
			lobbyId: l.lobbyId,
			ranked:  !l.abandoned,
			players: append([]gametype.Player(nil), g.Players...),
			seats:   seats,
		})
	}
}

func (l *Lobby) allReady() bool {
	for _, p := range l.game.Players {
		if !l.ready[p.ID] {
			return false
		}
	}
	return true
}

// deadline replaces the current phase timer
func (l *Lobby) deadline(d time.Duration) {
	if l.timer != nil {
		l.timer.Stop()
	}
	l.tick++
	l.timer = l.after(d, &phaseTimeout{tick: l.tick})
}

// timers run on their own goroutine, so they report back via a message
// instead of touching lobby state directly
func (l *Lobby) after(d time.Duration, msg any) *time.Timer {
	engine, self := l.ctx.Engine(), l.ctx.PID()
	return time.AfterFunc(d, func() { engine.Send(self, msg) })
}

func (l *Lobby) stopTimers() {
	if l.timer != nil {
		l.timer.Stop()
	}
	if l.startTimer != nil {
		l.startTimer.Stop()
	}
	l.tick++
}

// tell every player in the lobby who's in and when the game starts
func (l *Lobby) broadcastLobby() {
	players := make([]gametype.Player, len(l.seats))
	for i, s := range l.seats {
		players[i] = s.player
	}
	startsIn := 0
	if l.startTimer != nil {
		startsIn = int(math.Ceil(time.Until(l.startsAt).Seconds()))
	}
	for _, s := range l.seats {
		l.ctx.Send(s.pid, &types.LobbyUpdate{
			LobbyID:    l.lobbyId,
			You:        s.player.ID,
			Players:    players,
			MaxPlayers: maxLobbyPlayers,
			StartsIn:   startsIn,
		})
	}
}

func (l *Lobby) broadcastGame() {
	view := l.game.View()
	for _, s := range l.seats {
		l.ctx.Send(s.pid, &types.GameUpdate{You: s.player.ID, View: view})
	}
}

// the guesser already knows where their own dial is
func (l *Lobby) broadcastDial(guesser *seat) {
	update := &types.DialUpdate{PlayerID: guesser.player.ID, Value: l.game.Dial}
	for _, s := range l.seats {
		if s != guesser {
			l.ctx.Send(s.pid, update)
		}
	}
}

// "Priya Shah" -> "PS", "Sam" -> "SA", "Player 3" -> "P3"
func initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) >= 2 {
		return strings.ToUpper(string([]rune(parts[0])[:1]) + string([]rune(parts[1])[:1]))
	}
	r := []rune(name)
	return strings.ToUpper(string(r[:min(2, len(r))]))
}
