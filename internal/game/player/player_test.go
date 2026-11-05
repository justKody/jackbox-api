package player

import "testing"

func TestInitials(t *testing.T) {
	cases := map[string]string{
		"Priya Shah":      "PS",
		"Sam":             "SA",
		"Player 3":        "P3",
		"jo":              "JO",
		"A":               "A",
		"ana maria lopez": "AM",
		"Émile Zola":      "ÉZ",
	}
	for name, want := range cases {
		if got := Initials(name); got != want {
			t.Errorf("Initials(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestSeatColor(t *testing.T) {
	if SeatColor(1) != seatColors[0] {
		t.Fatalf("first seat should get the first color")
	}
	if SeatColor(len(seatColors)+1) != SeatColor(1) {
		t.Fatalf("colors should repeat past the last seat")
	}
}
