package player

import (
	"fmt"
	"strings"
)

// seat colors, in join order
var seatColors = []string{"#FFD23F", "#4C8DFF", "#FF8A3D", "#B58CFF", "#5ED3A8"}

// SeatColor is the color of the n-th seat handed out (1 based), colors repeat past the last one
func SeatColor(n int) string {
	return seatColors[(n-1)%len(seatColors)]
}

// DefaultName names a player who didn't pick one by the seat they got (1 based)
func DefaultName(n int) string {
	return fmt.Sprintf("Player %d", n)
}

// Initials of a name: "Priya Shah" -> "PS", "Sam" -> "SA", "Player 3" -> "P3"
func Initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) >= 2 {
		return strings.ToUpper(string([]rune(parts[0])[:1]) + string([]rune(parts[1])[:1]))
	}
	r := []rune(name)
	return strings.ToUpper(string(r[:min(2, len(r))]))
}
