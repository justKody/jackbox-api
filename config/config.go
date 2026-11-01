package config

import (
	"log"
	"os"
	"strings"

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
		Port:         getEnv("PORT", "4000"),
		AssetDir:     getEnv("ASSET_DIR", "s3asset"),
		AssetBaseURL: strings.TrimSuffix(getEnv("ASSET_BASE_URL", "/assets"), "/"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
