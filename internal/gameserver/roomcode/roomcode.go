package roomcode

import "math/rand"

const (
	Length = 4
	chars  = "ABCDEFGHJKLMNPQRSTUVWXYZ" // no I/O so codes are easy to read out loud
)

// New returns a jackbox style room code, e.g. "QXRT", that taken reports as free
func New(taken func(code string) bool) string {
	for {
		b := make([]byte, Length)
		for i := range b {
			b[i] = chars[rand.Intn(len(chars))]
		}
		if code := string(b); !taken(code) {
			return code
		}
	}
}
