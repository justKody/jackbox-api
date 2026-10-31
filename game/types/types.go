package types

import "github.com/justKody/jackbox-api/game/question"

type Call string

const (
	Higher Call = "higher"
	Lower  Call = "lower"
)

func (c Call) Valid() bool {
	return c == Higher || c == Lower
}

type Player struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Initials string `json:"initials"`
	Color    string `json:"color"`
	Score    int    `json:"score"`
	// the player's character, the game never reads these, they're passed on so
	// every client can draw everyone's animal
	Animal string `json:"animal"`
	Skin   string `json:"skin"`
}

type Phase string

const (
	PhaseGuess    Phase = "guess"
	PhaseVote     Phase = "vote"
	PhaseReveal   Phase = "reveal"
	PhaseGameOver Phase = "gameOver"
)

type RoundResult struct {
	PlayerID string `json:"playerId"`
	// what this player did this round, for the results list
	Label   string `json:"label"`
	Points  int    `json:"points"`
	Correct bool   `json:"correct"`
}

type State struct {
	Phase      Phase               `json:"phase"`
	Players    []Player            `json:"players"`
	Questions  []question.Question `json:"questions"`
	Round      int                 `json:"round"` // 0 based
	GuesserIdx int                 `json:"guesserIndex"`
	Guess      *int                `json:"guess"` // the guesser's locked guess, nil until locked
	Votes      map[string]Call     `json:"votes"` // player id -> call, guesser never votes
	Results    []RoundResult       `json:"results"`
}
