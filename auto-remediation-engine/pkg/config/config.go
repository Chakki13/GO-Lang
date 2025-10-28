package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the auto-remediation engine
// WHY: Centralize config so it's easy to manage via environment variables
type Config struct {
	// Teams webhook for alerts
	TeamsWebhookURL string

	// Git configuration for committing mutations
	GitRepoURL    string
	GitBranch     string
	GitAuthor     string
	GitEmail      string

	// Cluster info
	ClusterName   string
	ClusterRegion string

	// Feature flags - which remediations are enabled?
	EnablePodRestart      bool
	EnableMemoryIncrease  bool
	EnableSecurityMutations bool
	EnableProbeInjection  bool
	EnableCertAlerts      bool

	// Remediation parameters
	MemoryIncrementPercent int // How much to increase memory by
	CertExpiryThresholdDays int // Alert when cert expires in X days

	// Logging
	LogLevel string
}

// LoadConfig loads configuration from environment variables
// WHY: Environment variables allow container orchestration (K8s ConfigMaps)
// to control behavior without rebuilding the image
func LoadConfig() *Config {
	return &Config{
		TeamsWebhookURL:         getEnv("TEAMS_WEBHOOK_URL", ""),
		GitRepoURL:              getEnv("GIT_REPO_URL", ""),
		GitBranch:               getEnv("GIT_BRANCH", "main"),
		GitAuthor:               getEnv("GIT_AUTHOR", "auto-remediation"),
		GitEmail:                getEnv("GIT_EMAIL", "auto@remediation.local"),
		ClusterName:             getEnv("CLUSTER_NAME", "unknown"),
		ClusterRegion:           getEnv("CLUSTER_REGION", "us-east-1"),
		EnablePodRestart:        getEnvBool("ENABLE_POD_RESTART", true),
		EnableMemoryIncrease:    getEnvBool("ENABLE_MEMORY_INCREASE", true),
		EnableSecurityMutations: getEnvBool("ENABLE_SECURITY_MUTATIONS", true),
		EnableProbeInjection:    getEnvBool("ENABLE_PROBE_INJECTION", false), // Manual review needed
		EnableCertAlerts:        getEnvBool("ENABLE_CERT_ALERTS", true),
		MemoryIncrementPercent:  getEnvInt("MEMORY_INCREMENT_PERCENT", 50),
		CertExpiryThresholdDays: getEnvInt("CERT_EXPIRY_THRESHOLD_DAYS", 7),
		LogLevel:                getEnv("LOG_LEVEL", "info"),
	}
}

// Helper functions to get env vars with defaults
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}

func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}
