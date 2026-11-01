package config

import (
	"os"
	"strings"
)

// EnvOr is the value of the environment variable key, or fallback when it's unset or empty
func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Addr is where the game server listens, set PORT to change it
func Addr() string {
	return ":" + EnvOr("PORT", "4000")
}

// AssetDir is the s3asset folder the game server serves at /assets/, relative to
// the working dir (make server runs from backend/), set ASSET_DIR to change it
func AssetDir() string {
	return EnvOr("ASSET_DIR", "s3asset")
}

// AssetBaseURL is where clients fetch assets from, without a trailing slash.
// Set ASSET_BASE_URL to point at a bucket or CDN instead of the game server
func AssetBaseURL() string {
	return strings.TrimSuffix(EnvOr("ASSET_BASE_URL", "/assets"), "/")
}
