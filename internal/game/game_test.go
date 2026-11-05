package game

import (
	gametype "github.com/justKody/jackbox-api/internal/game/types"
	"math/rand/v2"
	"testing"
)

func newTestGame(t *testing.T, ids ...string) *Game {
	t.Helper()
	players := make([]gametype.Player, len(ids))
	for i, id := range ids {
		players[i] = gametype.Player{ID: id, Name: id}
	}
	g, err := NewGame(players, rand.New(rand.NewPCG(1, 2)))
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func TestScoreRound(t *testing.T) {
	players := []gametype.Player{{ID: "a"}, {ID: "b"}, {ID: "c"}, {ID: "d"}}
	votes := map[string]gametype.Call{"b": gametype.Higher, "c": gametype.Lower}
	r := ScoreRound(players, "a", 40, 50, votes)

	want := []int{90, CorrectCallPoints, 0, 0}
	for i, w := range want {
		if r[i].Points != w {
			t.Errorf("player %s: got %d points, want %d", r[i].PlayerID, r[i].Points, w)
		}
	}
	if r[3].Label != "No call" {
		t.Errorf("got label %q, want No call", r[3].Label)
	}

	// a bullseye gets the bonus and nobody wins the call
	r = ScoreRound(players, "a", 50, 50, votes)
	if r[0].Points != ExactGuessPoints || r[1].Points != 0 || r[2].Points != 0 {
		t.Errorf("bullseye scored wrong: %+v", r)
	}
}

func TestFullGame(t *testing.T) {
	g := newTestGame(t, "a", "b", "c")

	for round := 0; round < g.TotalRounds(); round++ {
		guesser := g.Guesser()
		if want := g.Players[round%3].ID; guesser.ID != want {
			t.Fatalf("round %d: guesser %s, want %s", round, guesser.ID, want)
		}
		if err := g.Vote("b", gametype.Higher); err != ErrWrongPhase {
			t.Fatalf("voting before the guess: got %v", err)
		}
		if g.View().Question.Answer != nil {
			t.Fatal("answer leaked before the reveal")
		}
		if err := g.LockGuess(guesser.ID, 50); err != nil {
			t.Fatal(err)
		}
		if err := g.Vote(guesser.ID, gametype.Higher); err != ErrGuesserVoting {
			t.Fatalf("guesser voting: got %v", err)
		}
		for _, p := range g.Voters() {
			if err := g.Vote(p.ID, gametype.Higher); err != nil {
				t.Fatal(err)
			}
		}
		if !g.AllVoted() {
			t.Fatal("expected everyone voted")
		}
		if err := g.Reveal(); err != nil {
			t.Fatal(err)
		}
		if g.View().Question.Answer == nil {
			t.Fatal("answer hidden after the reveal")
		}
		if err := g.Next(); err != nil {
			t.Fatal(err)
		}
	}
	if g.Phase != gametype.PhaseGameOver {
		t.Fatalf("phase %s, want gameOver", g.Phase)
	}
}

func TestRemoveGuesser(t *testing.T) {
	g := newTestGame(t, "a", "b", "c")
	if err := g.LockGuess("a", 30); err != nil {
		t.Fatal(err)
	}
	if !g.RemovePlayer("a") {
		t.Fatal("game should go on with two players")
	}
	if g.Phase != gametype.PhaseGuess || g.Guesser().ID != "b" {
		t.Fatalf("got phase %s guesser %s, want guess phase with b", g.Phase, g.Guesser().ID)
	}
	if g.RemovePlayer("b") {
		t.Fatal("game can't go on with one player")
	}
}

func TestRemoveLastGuesser(t *testing.T) {
	g := newTestGame(t, "a", "b")
	g.LockGuess("a", 50)
	g.Vote("b", gametype.Higher)
	g.Reveal()
	g.Next() // b guesses round 2
	if g.RemovePlayer("b") {
		t.Fatal("game can't go on with one player")
	}
	// the view still has to work for the player left at the table
	if v := g.View(); v.GuesserID != "a" {
		t.Fatalf("got guesser %q, want a", v.GuesserID)
	}
}

func TestDial(t *testing.T) {
	g := newTestGame(t, "a", "b")
	if err := g.MoveDial("b", 70); err != ErrNotGuesser {
		t.Fatalf("only the guesser turns the dial, got %v", err)
	}
	if err := g.MoveDial("a", 140); err != nil || g.View().Dial != 100 {
		t.Fatalf("dial should clamp to 100, got %d err %v", g.View().Dial, err)
	}
	g.MoveDial("a", 72)
	// running out of time locks in wherever the dial was left
	if err := g.GuessTimeout(); err != nil || *g.Guess != 72 {
		t.Fatalf("timeout should lock 72, got %v err %v", g.Guess, err)
	}
	if err := g.MoveDial("a", 10); err != ErrWrongPhase {
		t.Fatalf("dial is locked once the guess is in, got %v", err)
	}
}

// every round's question goes out with the Hat's reading of it
func TestViewSendsVoiceOver(t *testing.T) {
	g := newTestGame(t, "a", "b", "c")
	for round := 0; round < g.TotalRounds(); round++ {
		g.Round = round
		q := g.View().Question
		if q.VoiceOver != "/assets/v2/"+q.ID+".mp3" {
			t.Errorf("round %d: question %s sent voice-over %q", round, q.ID, q.VoiceOver)
		}
	}
}
