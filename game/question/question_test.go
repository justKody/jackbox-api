package question

import (
	"os"
	"path/filepath"
	"testing"
)

// every question the game can ask comes with the Hat reading it
func TestBankHasVoiceOvers(t *testing.T) {
	if len(Bank) == 0 {
		t.Fatal("empty question bank, run build_dataset.py")
	}
	for _, q := range Bank {
		if q.VoiceOver == "" {
			t.Errorf("%s has no voice-over", q.ID)
			continue
		}
		if _, err := os.Stat(filepath.Join("..", "..", "s3asset", q.VoiceOver)); err != nil {
			t.Errorf("%s: %v", q.ID, err)
		}
		if q.Answer < 0 || q.Answer > 100 {
			t.Errorf("%s: answer %d out of range", q.ID, q.Answer)
		}
	}
}

func TestVoiceURL(t *testing.T) {
	q := Question{VoiceOver: "v2/owid1.mp3"}
	if got := q.VoiceURL(); got != "/assets/v2/owid1.mp3" {
		t.Errorf("VoiceURL() = %q", got)
	}
	if got := (Question{}).VoiceURL(); got != "" {
		t.Errorf("no voice-over should give no url, got %q", got)
	}
}
