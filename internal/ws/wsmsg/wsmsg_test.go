package wsmsg

import (
	"encoding/json"
	"testing"

	"github.com/justKody/jackbox-api/internal/ws/types"
)

func TestDecode(t *testing.T) {
	var m types.LockGuessMsg
	if err := Decode(types.WSMessage{Type: "lockGuess", Data: json.RawMessage(`{"guess":42}`)}, &m); err != nil || m.Guess != 42 {
		t.Fatalf("want guess 42, got %d err %v", m.Guess, err)
	}

	// no data is fine, v keeps its zero value
	var p types.Profile
	if err := Decode(types.WSMessage{Type: "findMatch"}, &p); err != nil || p != (types.Profile{}) {
		t.Fatalf("empty data should decode to the zero value, got %+v err %v", p, err)
	}

	if err := Decode(types.WSMessage{Type: "vote", Data: json.RawMessage(`{"call":`)}, &m); err == nil {
		t.Fatal("malformed data should be an error")
	}
}

func TestAllowOrigin(t *testing.T) {
	cases := []struct {
		origin, host string
		want         bool
	}{
		{"", "example.com", true},
		{"http://localhost:8081", "example.com", true},
		{"http://192.168.1.5:4000", "192.168.1.5:4000", true},
		{"http://evil.com", "example.com", false},
	}
	for _, c := range cases {
		if got := AllowOrigin(c.origin, c.host); got != c.want {
			t.Errorf("AllowOrigin(%q, %q) = %v, want %v", c.origin, c.host, got, c.want)
		}
	}
}
