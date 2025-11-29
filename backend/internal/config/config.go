package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	// Server
	Port        string
	Environment string
	GinMode     string

	// Database
	DatabaseURL       string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime int // seconds

	// Security
	JWTSecret      string
	APIKeySecret   string
	AllowedOrigins []string

	// CORS
	CORSAllowOrigins     []string
	CORSAllowCredentials bool

	// Rate Limiting
	RateLimitEnabled bool
	RateLimitRPS     int
	RateLimitBurst   int

	// Logging
	LogLevel  string
	LogFormat string
	LogFile   string

	// External Integrations
	ConvinAPIKey           string
	ConvinAPIURL           string
	ConvinWebhookSecret    string
	TwilioAccountSID       string
	TwilioAuthToken        string
	TwilioWebhookSecret    string
	TelephonyWebhookSecret string

	// Monitoring
	SentryDSN          string
	NewRelicLicenseKey string

	// Feature Flags
	EnableWebhooks           bool
	EnableRealtimeProcessing bool
	EnableFraudDetection     bool
}

func Load() *Config {
	env := getEnv("ENVIRONMENT", "development")

	return &Config{
		// Server
		Port:        getEnv("PORT", "8080"),
		Environment: env,
		GinMode:     getEnv("GIN_MODE", getGinMode(env)),

		// Database
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://localhost/convin_crae?sslmode=disable"),
		DBMaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 300),

		// Security
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		APIKeySecret:   getEnv("API_KEY_SECRET", "change-me-in-production"),
		AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),

		// CORS
		CORSAllowOrigins:     getEnvAsSlice("CORS_ALLOW_ORIGINS", []string{"http://localhost:3000"}),
		CORSAllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),

		// Rate Limiting
		RateLimitEnabled: getEnvAsBool("RATE_LIMIT_ENABLED", true),
		RateLimitRPS:     getEnvAsInt("RATE_LIMIT_RPS", 100),
		RateLimitBurst:   getEnvAsInt("RATE_LIMIT_BURST", 200),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", getLogLevel(env)),
		LogFormat: getEnv("LOG_FORMAT", getLogFormat(env)),
		LogFile:   getEnv("LOG_FILE", "logs/app.log"),

		// External Integrations
		ConvinAPIKey:           getEnv("CONVIN_API_KEY", ""),
		ConvinAPIURL:           getEnv("CONVIN_API_URL", "https://api.convin.ai"),
		ConvinWebhookSecret:    getEnv("CONVIN_WEBHOOK_SECRET", ""),
		TwilioAccountSID:       getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:        getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioWebhookSecret:    getEnv("TWILIO_WEBHOOK_SECRET", ""),
		TelephonyWebhookSecret: getEnv("TELEPHONY_WEBHOOK_SECRET", ""),

		// Monitoring
		SentryDSN:          getEnv("SENTRY_DSN", ""),
		NewRelicLicenseKey: getEnv("NEW_RELIC_LICENSE_KEY", ""),

		// Feature Flags
		EnableWebhooks:           getEnvAsBool("ENABLE_WEBHOOKS", true),
		EnableRealtimeProcessing: getEnvAsBool("ENABLE_REALTIME_PROCESSING", true),
		EnableFraudDetection:     getEnvAsBool("ENABLE_FRAUD_DETECTION", true),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func getGinMode(env string) string {
	if env == "production" {
		return "release"
	}
	return "debug"
}

func getLogLevel(env string) string {
	if env == "production" {
		return "warn"
	}
	return "debug"
}

func getLogFormat(env string) string {
	if env == "production" {
		return "json"
	}
	return "text"
}
