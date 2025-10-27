package config

import "os"

// Config holds all configuration for the app
type Config struct {
	TeamsWebhookURL string
	ClusterName     string
	CostThreshold   float64 // Alert threshold in dollars per hour
	LogLevel        string
}

// LoadConfig loads config from environment variables
func LoadConfig() *Config {
	return &Config{
		TeamsWebhookURL: os.Getenv("TEAMS_WEBHOOK_URL"),
		ClusterName:     os.Getenv("CLUSTER_NAME"),
		CostThreshold:   50.0, // Default: alert if >$50/hr
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv gets an env var with a default
func getEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
