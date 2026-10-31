package game

import (
	"errors"
	"math/rand/v2"
	"time"

	"github.com/justKody/jackbox-api/game/question"
	gametype "github.com/justKody/jackbox-api/game/types"
)

const (
	Rounds       = 6
	GuessSeconds = 20 * time.Second
	VoteSeconds  = 15 * time.Second

	// what the dial starts at, locked in if the guesser runs out of time
	DefaultGuess = 50
)

var (
	ErrWrongPhase     = errors.New("action not allowed in this phase")
	ErrNotGuesser     = errors.New("only the guesser can do that")
	ErrGuesserVoting  = errors.New("the guesser can't vote")
	ErrAlreadyVoted   = errors.New("player already voted")
	ErrUnknownPlayer  = errors.New("player not in this game")
	ErrInvalidCall    = errors.New("call must be higher or lower")
	ErrNotEnoughSeats = errors.New("need at least two players")
)

// Game is the authoritative state of one match. It has no timers or locking,
// the lobby that owns it drives it (e.g. calling GuessTimeout after GuessSeconds)
type Game struct {
	gametype.State
	rng *rand.Rand
}

func NewGame(players []gametype.Player, rng *rand.Rand) (*Game, error) {
	if len(players) < 2 {
		return nil, ErrNotEnoughSeats
	}
	if rng == nil {
		rng = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}
	g := &Game{rng: rng}
	g.reset(players)
	return g, nil
}

func (g *Game) reset(players []gametype.Player) {
	ps := make([]gametype.Player, len(players))
	for i, p := range players {
		p.Score = 0
		ps[i] = p
	}
	qs := make([]question.Question, len(question.Bank))
	copy(qs, question.Bank)
	g.rng.Shuffle(len(qs), func(i, j int) { qs[i], qs[j] = qs[j], qs[i] })

	g.State = gametype.State{
		Players:   ps,
		Questions: qs[:min(Rounds, len(qs))],
	}
	g.beginRound()
}

// enter a round, the guesser gets the dial
func (g *Game) beginRound() {
	g.Phase = gametype.PhaseGuess
	g.Guess = nil
	g.Votes = make(map[string]gametype.Call)
	g.Results = nil
}

func (g *Game) Restart() {
	g.reset(g.Players)
}

func (g *Game) Guesser() gametype.Player {
	return g.Players[g.GuesserIdx]
}

func (g *Game) Question() question.Question {
	return g.Questions[g.Round]
}

func (g *Game) TotalRounds() int {
	return len(g.Questions)
}

func (g *Game) IsLastRound() bool {
	return g.Round+1 >= len(g.Questions)
}

// Voters is everyone except the guesser
func (g *Game) Voters() []gametype.Player {
	voters := make([]gametype.Player, 0, len(g.Players)-1)
	for i, p := range g.Players {
		if i != g.GuesserIdx {
			voters = append(voters, p)
		}
	}
	return voters
}

func (g *Game) AllVoted() bool {
	for _, p := range g.Voters() {
		if _, ok := g.Votes[p.ID]; !ok {
			return false
		}
	}
	return true
}

func (g *Game) LockGuess(playerId string, guess int) error {
	if g.Phase != gametype.PhaseGuess {
		return ErrWrongPhase
	}
	if g.Guesser().ID != playerId {
		return ErrNotGuesser
	}
	v := clamp(guess)
	g.Guess = &v
	g.Phase = gametype.PhaseVote
	return nil
}

// GuessTimeout locks in the dial's starting value when the guesser runs out of time
func (g *Game) GuessTimeout() error {
	return g.LockGuess(g.Guesser().ID, DefaultGuess)
}

func (g *Game) Vote(playerId string, call gametype.Call) error {
	if g.Phase != gametype.PhaseVote {
		return ErrWrongPhase
	}
	if !call.Valid() {
		return ErrInvalidCall
	}
	if g.playerIndex(playerId) < 0 {
		return ErrUnknownPlayer
	}
	if g.Guesser().ID == playerId {
		return ErrGuesserVoting
	}
	if _, ok := g.Votes[playerId]; ok {
		return ErrAlreadyVoted
	}
	g.Votes[playerId] = call
	return nil
}

// Reveal scores the round, call it once everyone voted or the vote timer ran out
func (g *Game) Reveal() error {
	if g.Phase != gametype.PhaseVote || g.Guess == nil {
		return ErrWrongPhase
	}
	results := ScoreRound(g.Players, g.Guesser().ID, *g.Guess, g.Question().Answer, g.Votes)
	for i := range g.Players {
		g.Players[i].Score += results[i].Points
	}
	g.Results = results
	g.Phase = gametype.PhaseReveal
	return nil
}

// Next moves on from the reveal, to the next round or to game over
func (g *Game) Next() error {
	if g.Phase != gametype.PhaseReveal {
		return ErrWrongPhase
	}
	if g.IsLastRound() {
		g.Phase = gametype.PhaseGameOver
		return nil
	}
	g.Round++
	g.GuesserIdx = (g.GuesserIdx + 1) % len(g.Players)
	g.beginRound()
	return nil
}

// RemovePlayer drops a player who left mid game, returns false if fewer than
// two players remain and the game can't go on
func (g *Game) RemovePlayer(playerId string) bool {
	i := g.playerIndex(playerId)
	if i < 0 {
		return len(g.Players) >= 2
	}
	wasGuesser := i == g.GuesserIdx
	g.Players = append(g.Players[:i], g.Players[i+1:]...)
	delete(g.Votes, playerId)
	if len(g.Players) < 2 {
		return false
	}
	if i < g.GuesserIdx {
		g.GuesserIdx--
	}
	g.GuesserIdx %= len(g.Players)
	// the round can't continue without its guesser, hand the dial to the next player
	if wasGuesser && (g.Phase == gametype.PhaseGuess || g.Phase == gametype.PhaseVote) {
		g.beginRound()
	}
	return true
}

func (g *Game) playerIndex(playerId string) int {
	for i, p := range g.Players {
		if p.ID == playerId {
			return i
		}
	}
	return -1
}

func clamp(n int) int {
	return max(0, min(100, n))
}
