package game

import (
	"fmt"

	gametype "github.com/justKody/jackbox-api/game/types"
)

const (
	CorrectCallPoints = 100
	ExactGuessPoints  = 150
)

// CorrectCall returns which call was right, ok is false when the guess is
// exactly the answer (nobody wins the call)
func CorrectCall(guess, answer int) (call gametype.Call, ok bool) {
	switch {
	case answer > guess:
		return gametype.Higher, true
	case answer < guess:
		return gametype.Lower, true
	}
	return "", false
}

// ClosenessPoints gives the guesser 100 minus how far off they were, a bullseye earns a bonus
func ClosenessPoints(guess, answer int) int {
	off := abs(guess - answer)
	if off == 0 {
		return ExactGuessPoints
	}
	return max(0, 100-off)
}

func ScoreRound(players []gametype.Player, guesserId string, guess, answer int, votes map[string]gametype.Call) []gametype.RoundResult {
	right, ok := CorrectCall(guess, answer)
	results := make([]gametype.RoundResult, len(players))
	for i, p := range players {
		if p.ID == guesserId {
			points := ClosenessPoints(guess, answer)
			results[i] = gametype.RoundResult{
				PlayerID: p.ID,
				Label:    fmt.Sprintf("Guessed %d%% · closeness", guess),
				Points:   points,
				Correct:  points > 0,
			}
			continue
		}
		call, voted := votes[p.ID]
		if !voted {
			results[i] = gametype.RoundResult{PlayerID: p.ID, Label: "No call"}
			continue
		}
		correct := ok && call == right
		points := 0
		if correct {
			points = CorrectCallPoints
		}
		results[i] = gametype.RoundResult{
			PlayerID: p.ID,
			Label:    fmt.Sprintf("Called %s", call),
			Points:   points,
			Correct:  correct,
		}
	}
	return results
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
