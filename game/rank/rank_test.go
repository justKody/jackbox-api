package rank

import (
	"testing"

	gametype "github.com/justKody/jackbox-api/game/types"
)

func TestRanked(t *testing.T) {
	if Delta(1, 5) != 30 || Delta(5, 5) != -20 || Delta(3, 5) != 5 {
		t.Error("rp delta off")
	}
	if got := For(1150).Label; got != "Conjurer II" {
		t.Errorf("got %q, want Conjurer II", got)
	}
	if got := For(1400).Label; got != "Archmage" {
		t.Errorf("got %q, want Archmage", got)
	}

	players := []gametype.Player{{ID: "a", Score: 300}, {ID: "b", Score: 300}, {ID: "c", Score: 100}}
	ratings := map[string]int{"a": 1010, "c": 5}
	res := Apply(players, ratings)
	// a and b tie for first
	if res[0].Delta != 30 || res[1].Place != 1 || res[1].After != StartingRP+30 {
		t.Errorf("tie handled wrong: %+v", res)
	}
	if ratings["c"] != 0 {
		t.Errorf("rating went below zero: %d", ratings["c"])
	}
}
