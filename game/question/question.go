package question

import (
	_ "embed"
	"encoding/json"

	"github.com/justKody/jackbox-api/config"
)

type Question struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	// the true answer, 0-100 (%)
	Answer int `json:"answer"`
	// the Hat reading the question aloud, a path under s3asset
	VoiceOver string `json:"voiceOver"`
}

// built by build_dataset.py: the OWID questions that have a voice-over
//
//go:embed dataset.json
var dataset []byte

// question bank, each answer is a percentage (0-100)
// the game picks up to game.Rounds of them at random each game
var Bank = load()

func load() []Question {
	var qs []Question
	if err := json.Unmarshal(dataset, &qs); err != nil {
		panic("question: bad dataset.json: " + err.Error())
	}
	return qs
}

// VoiceURL is where a client can fetch the voice-over, "" when there is none
func (q Question) VoiceURL() string {
	if q.VoiceOver == "" {
		return ""
	}
	return config.Envs.AssetBaseURL + "/" + q.VoiceOver
}
