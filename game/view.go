package game

import (
	"maps"
	"math"
	"slices"

	gametype "github.com/justKody/jackbox-api/game/types"
)

// QuestionView is a question as clients see it, the answer stays hidden until the reveal
type QuestionView struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	// url of the Hat reading the question, played as the question appears
	VoiceOver string `json:"voiceOver,omitempty"`
	Answer    *int   `json:"answer,omitempty"`
}

// View is what gets sent to clients. State itself carries every answer and
// every call, so it should never go over the wire as is
type View struct {
	Phase       gametype.Phase    `json:"phase"`
	Players     []gametype.Player `json:"players"`
	Round       int               `json:"round"` // 0 based
	TotalRounds int               `json:"totalRounds"`
	GuesserID   string            `json:"guesserId"`
	Question    QuestionView      `json:"question"`
	Guess       *int              `json:"guess"`
	Dial        int               `json:"dial"` // the guesser's dial, live while they're guessing
	// how long each timed phase lasts, so clients count down the same clock the server keeps
	GuessSeconds int `json:"guessSeconds"`
	VoteSeconds  int `json:"voteSeconds"`
	// this round's reveal and standings together, the next round starts by itself after it.
	// it can differ round to round, so clients read it from every update
	RevealSeconds int `json:"revealSeconds"`
	// who has locked in a call, so the table can show "Locked" without giving calls away
	Locked  []string                 `json:"locked"`
	Votes   map[string]gametype.Call `json:"votes,omitempty"` // only once revealed
	Results []gametype.RoundResult   `json:"results,omitempty"`
}

func (g *Game) View() View {
	q := g.Question()
	revealed := g.Phase == gametype.PhaseReveal || g.Phase == gametype.PhaseGameOver

	locked := make([]string, 0, len(g.Votes))
	for _, p := range g.Players {
		if _, ok := g.Votes[p.ID]; ok {
			locked = append(locked, p.ID)
		}
	}

	v := View{
		Phase:       g.Phase,
		Players:     slices.Clone(g.Players),
		Round:       g.Round,
		TotalRounds: g.TotalRounds(),
		GuesserID:   g.Guesser().ID,
		Question:    QuestionView{ID: q.ID, Text: q.Text, VoiceOver: q.VoiceURL()},
		Guess:       g.Guess,
		Dial:        g.Dial,
		Locked:      locked,

		GuessSeconds: int(GuessSeconds.Seconds()),
		VoteSeconds:  int(VoteSeconds.Seconds()),
	}
	if revealed {
		// rounded up, so a client never moves on before the server does
		v.RevealSeconds = int(math.Ceil(g.RevealFor.Seconds()))
		answer := q.Answer
		v.Question.Answer = &answer
		// copied, the view is read on other goroutines while the game keeps going
		v.Votes = maps.Clone(g.Votes)
		v.Results = slices.Clone(g.Results)
	}
	return v
}
