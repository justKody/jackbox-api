package rank

import (
	"math"

	gametype "github.com/justKody/jackbox-api/game/types"
)

// the ranked ladder: every player has a rating (RP)
// only ranked games move RP, lobby games are just for fun

const StartingRP = 1000

type Tier struct {
	Name  string `json:"name"`
	Min   int    `json:"min"`
	Color string `json:"color"`
}

// lowest first, each tier spans 100 RP and is split into divisions III -> I
var Tiers = []Tier{
	{Name: "Apprentice", Min: 0, Color: "#A9B1C9"},
	{Name: "Adept", Min: 1000, Color: "#5CE1FF"},
	{Name: "Conjurer", Min: 1100, Color: "#9ACD5E"},
	{Name: "Sorcerer", Min: 1200, Color: "#A98BFF"},
	{Name: "Archmage", Min: 1300, Color: "#F4C15D"},
}

type Info struct {
	Tier Tier `json:"tier"`
	// e.g. "II", empty for the top tier which has no divisions
	Division string `json:"division"`
	Label    string `json:"label"`
	// 0-1 progress through the current tier
	Progress float64 `json:"progress"`
	Next     *Tier   `json:"next,omitempty"`
	ToNext   int     `json:"toNext"`
}

func For(rp int) Info {
	i := len(Tiers) - 1
	for i > 0 && rp < Tiers[i].Min {
		i--
	}
	tier := Tiers[i]
	into := float64(max(0, rp-tier.Min))

	if i == len(Tiers)-1 {
		return Info{Tier: tier, Label: tier.Name, Progress: math.Min(1, into/300)}
	}

	next := Tiers[i+1]
	progress := math.Min(1, into/float64(next.Min-tier.Min))
	division := []string{"III", "II", "I"}[min(2, int(math.Floor(progress*3)))]
	return Info{
		Tier:     tier,
		Division: division,
		Label:    tier.Name + " " + division,
		Progress: progress,
		Next:     &next,
		ToNext:   next.Min - rp,
	}
}

// Delta is the RP for finishing place (1 based) out of of: winner +30, last -20, linear between
func Delta(place, of int) int {
	if of < 2 {
		return 0
	}
	return int(math.Round(30 - float64(50*(place-1))/float64(of-1)))
}

// Place is the player's 1 based finishing place, ties share the better place
func Place(players []gametype.Player, p gametype.Player) int {
	place := 1
	for _, o := range players {
		if o.Score > p.Score {
			place++
		}
	}
	return place
}

type Result struct {
	PlayerID string `json:"playerId"`
	Before   int    `json:"before"`
	After    int    `json:"after"`
	Delta    int    `json:"delta"`
	Place    int    `json:"place"`
	Of       int    `json:"of"`
}

// Apply applies a finished ranked match to everyone at the table.
// ratings is player id -> current RP (missing players start at StartingRP),
// it is updated in place and each player's change is returned
func Apply(players []gametype.Player, ratings map[string]int) []Result {
	results := make([]Result, len(players))
	for i, p := range players {
		before, ok := ratings[p.ID]
		if !ok {
			before = StartingRP
		}
		place := Place(players, p)
		after := max(0, before+Delta(place, len(players)))
		ratings[p.ID] = after
		results[i] = Result{
			PlayerID: p.ID,
			Before:   before,
			After:    after,
			Delta:    after - before,
			Place:    place,
			Of:       len(players),
		}
	}
	return results
}
