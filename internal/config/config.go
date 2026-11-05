package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// port the game server listens on
	Port string
	// @this will be removed when indtroduction of s3
	// the s3asset folder the game server serves at /assets/, relative to the
	// working dir (make server runs from backend/)
	AssetDir string
	// where clients fetch assets from, without a trailing slash. Point it at a
	// bucket or CDN instead of the game server
	AssetBaseURL string
	// how long the guesser has to lock in, and the others to make their call.
	// clients are sent these with every game update and count down the same clock
	GuessSeconds time.Duration
	VoteSeconds  time.Duration
	// how long the answer and standings stay up before the next round starts on its own
	RevealSeconds time.Duration
}

// Envs is read once at startup, before any package level var that uses it
// (e.g. question.VoiceURL), since every package using it imports config
var Envs = initConfig()

func initConfig() Config {
	// Optional locally; in Docker/K8s env vars are injected directly.
	// vars already set in the environment win over .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return Config{
		Port:          getEnv("PORT", "4000"),
		AssetDir:      getEnv("ASSET_DIR", "s3asset"),
		AssetBaseURL:  strings.TrimSuffix(getEnv("ASSET_BASE_URL", "/assets"), "/"),
		GuessSeconds:  getEnvSeconds("GUESS_SECONDS", 10),
		VoteSeconds:   getEnvSeconds("VOTE_SECONDS", 10),
		RevealSeconds: getEnvSeconds("REVEAL_SECONDS", 15),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvSeconds reads a whole number of seconds, a bad value stops the server
// rather than run games on a clock nobody set
func getEnvSeconds(key string, fallback int) time.Duration {
	n := fallback
	if value, ok := os.LookupEnv(key); ok {
		var err error
		if n, err = strconv.Atoi(value); err != nil || n <= 0 {
			log.Fatalf("%s must be a positive whole number of seconds, got %q", key, value)
		}
	}
	return time.Duration(n) * time.Second
}
