package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	DatabaseURL string
	JWT         JWTConfig
	Server      ServerConfig
	OIDC        OIDCConfig
}

type JWTConfig struct {
	Secret string
	Issuer string
	TTL    time.Duration
}

type ServerConfig struct {
	Addr string
}

type OIDCConfig struct {
	Google GoogleOIDCConfig
}

type GoogleOIDCConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Issuer       string
}

func Load() *Config {
	return &Config{
		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
		JWT: JWTConfig{
			Secret: getEnvOrDefault("JWT_SECRET", "dev-secret-change-me"),
			Issuer: getEnvOrDefault("JWT_ISSUER", "aggipay"),
			TTL:    parseTTLHours(os.Getenv("JWT_TTL_HOURS"), 24*time.Hour),
		},
		Server: ServerConfig{
			Addr: getEnvOrDefault("SERVER_ADDR", ":8000"),
		},
		OIDC: OIDCConfig{
			Google: GoogleOIDCConfig{
				ClientID:     strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
				ClientSecret: strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET")),
				RedirectURL:  getEnvOrDefault("GOOGLE_REDIRECT_URL", "http://localhost:8000/auth/google/callback"),
				Issuer:       getEnvOrDefault("GOOGLE_OIDC_ISSUER", "https://accounts.google.com"),
			},
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return defaultValue
}

func parseTTLHours(v string, defaultTTL time.Duration) time.Duration {
	v = strings.TrimSpace(v)
	if v == "" {
		return defaultTTL
	}
	if parsed, err := time.ParseDuration(v + "h"); err == nil {
		return parsed
	}
	return defaultTTL
}
