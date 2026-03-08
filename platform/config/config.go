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
	Frontend    FrontendConfig
	OIDC        OIDCConfig
	Pepper      string
	RabbitMQ    RabbitMQConfig
	Asaas       AsaasConfig
	Clicksign   ClicksignConfig
	Webhooks    WebhooksConfig
}

type AsaasConfig struct {
	BaseURL string
	APIKey  string
}

type ClicksignConfig struct {
	BaseURL string
	Token   string
}

type WebhooksConfig struct {
	AsaasSecret     string
	ClicksignSecret string
}

type JWTConfig struct {
	Secret string
	Issuer string
	TTL    time.Duration
}

type ServerConfig struct {
	Addr string
}

type FrontendConfig struct {
	NextInternalURL string
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

type RabbitMQConfig struct {
	URL           string
	Exchange      string
	CommandsQueue string
	EventsQueue   string
	DLQQueue      string
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
		Frontend: FrontendConfig{
			NextInternalURL: getEnvOrDefault("NEXT_INTERNAL_URL", "http://127.0.0.1:3001"),
		},
		OIDC: OIDCConfig{
			Google: GoogleOIDCConfig{
				ClientID:     strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
				ClientSecret: strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET")),
				RedirectURL:  getEnvOrDefault("GOOGLE_REDIRECT_URL", "http://localhost:8000/auth/google/callback"),
				Issuer:       getEnvOrDefault("GOOGLE_OIDC_ISSUER", "https://accounts.google.com"),
			},
		},
		Pepper: getEnvOrDefault("PASSWORD_PEPPER", "dev-pepper-change-me"),
		RabbitMQ: RabbitMQConfig{
			URL:           getEnvOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Exchange:      getEnvOrDefault("RABBITMQ_EXCHANGE", "aggipay.payment"),
			CommandsQueue: getEnvOrDefault("RABBITMQ_COMMANDS_QUEUE", "payment.commands"),
			EventsQueue:   getEnvOrDefault("RABBITMQ_EVENTS_QUEUE", "payment.events"),
			DLQQueue:      getEnvOrDefault("RABBITMQ_DLQ_QUEUE", "payment.dlq"),
		},
		Asaas: AsaasConfig{
			BaseURL: getEnvOrDefault("ASAAS_BASE_URL", "https://api.asaas.com/v3"),
			APIKey:  strings.TrimSpace(os.Getenv("ASAAS_API_KEY")),
		},
		Clicksign: ClicksignConfig{
			BaseURL: getEnvOrDefault("CLICKSIGN_BASE_URL", "https://sandbox.clicksign.com/api/v1"),
			Token:   strings.TrimSpace(os.Getenv("CLICKSIGN_TOKEN")),
		},
		Webhooks: WebhooksConfig{
			AsaasSecret:     strings.TrimSpace(os.Getenv("ASAAS_WEBHOOK_SECRET")),
			ClicksignSecret: strings.TrimSpace(os.Getenv("CLICKSIGN_WEBHOOK_SECRET")),
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
