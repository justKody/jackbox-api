package game

import gametype "github.com/justKody/jackbox-api/game/types"

// QuestionView is a question as clients see it, the answer stays hidden until the reveal
type QuestionView struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Answer *int   `json:"answer,omitempty"`
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
		Players:     append([]gametype.Player(nil), g.Players...),
		Round:       g.Round,
		TotalRounds: g.TotalRounds(),
		GuesserID:   g.Guesser().ID,
		Question:    QuestionView{ID: q.ID, Text: q.Text},
		Guess:       g.Guess,
		Locked:      locked,
	}
	if revealed {
		answer := q.Answer
		v.Question.Answer = &answer
		// copied, the view is read on other goroutines while the game keeps going
		v.Votes = make(map[string]gametype.Call, len(g.Votes))
		for id, c := range g.Votes {
			v.Votes[id] = c
		}
		v.Results = append([]gametype.RoundResult(nil), g.Results...)
	}
	return v
}
