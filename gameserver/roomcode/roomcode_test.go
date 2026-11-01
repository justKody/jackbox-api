package roomcode

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	for range 100 {
		code := New(func(string) bool { return false })
		if len(code) != Length {
			t.Fatalf("want %d chars, got %q", Length, code)
		}
		if strings.ContainsAny(code, "IO") {
			t.Fatalf("code %q has an I or O", code)
		}
	}
}

func TestNewSkipsTaken(t *testing.T) {
	seen := map[string]bool{}
	tries := 0
	// the first few codes are taken, New has to keep going
	code := New(func(c string) bool {
		tries++
		seen[c] = true
		return tries <= 3
	})
	if tries != 4 || !seen[code] {
		t.Fatalf("want the 4th code, got %q after %d tries", code, tries)
	}
}
