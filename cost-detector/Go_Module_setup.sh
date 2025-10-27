#!/bin/bash

# Go Module Setup for Cost Detector
# This script initializes Go modules and creates placeholder files in each package
# Run this AFTER running scaffold.sh

set -e

# Colors for output
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Go Module Setup for Cost Detector${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check if we're in the cost-detector directory
if [ ! -d "cmd/cost-detector" ]; then
    echo "❌ Error: Must be in cost-detector directory"
    echo "Run: cd cost-detector"
    exit 1
fi

echo -e "${YELLOW}📝 What's happening:${NC}"
echo "1. Initializing Go module (go.mod)"
echo "2. Creating placeholder Go files in each package"
echo "3. Each file will have a basic structure to build on"
echo ""

# Step 1: Initialize Go module
echo -e "${GREEN}✅ Step 1: Initialize Go Module${NC}"
echo "Command: go mod init github.com/yourname/cost-detector"
echo ""
echo "What this does:"
echo "  - Creates go.mod file (like package.json for Node)"
echo "  - Tracks all dependencies your project uses"
echo "  - Lets Go know this is a project named 'cost-detector'"
echo ""

go mod init github.com/yourusername/cost-detector 2>/dev/null || echo "⚠️  go.mod might already exist"

echo ""
echo -e "${GREEN}✅ Step 2: Create Package Files${NC}"
echo ""

# Create models package (data structures)
echo "Creating pkg/models/models.go..."
cat > pkg/models/models.go << 'EOF'
package models

// Pod represents a Kubernetes pod
type Pod struct {
	Name       string  // Pod name
	Namespace  string  // Namespace it's in
	CPU        float64 // CPU requested (cores)
	Memory     float64 // Memory requested (GB)
	CostPerHr  float64 // Calculated hourly cost
}

// CostAlert represents a cost alert to send to Teams
type CostAlert struct {
	Team       string
	Service    string
	CostPerHr  float64
	Message    string
	Severity   string // "info", "warning", "critical"
}

// NodePrice holds pricing info for a node type
type NodePrice struct {
	InstanceType string  // e.g., "t3.large"
	CostPerHour  float64 // Cost per hour in dollars
}
EOF
echo "  ✓ Created"

# Create watcher package (watches pod events)
echo "Creating pkg/watcher/watcher.go..."
cat > pkg/watcher/watcher.go << 'EOF'
package watcher

import "fmt"

// Watcher monitors pod creation and deletion events
type Watcher struct {
	ClusterName string
}

// NewWatcher creates a new pod watcher
func NewWatcher(clusterName string) *Watcher {
	return &Watcher{
		ClusterName: clusterName,
	}
}

// Start begins watching pod events
func (w *Watcher) Start() error {
	fmt.Printf("Starting watcher for cluster: %s\n", w.ClusterName)
	// TODO: Connect to Kubernetes API
	// TODO: Watch for pod creation/deletion events
	return nil
}

// Stop stops watching pod events
func (w *Watcher) Stop() {
	fmt.Println("Stopping watcher")
	// TODO: Clean up connections
}
EOF
echo "  ✓ Created"

# Create calculator package (does cost math)
echo "Creating pkg/calculator/calculator.go..."
cat > pkg/calculator/calculator.go << 'EOF'
package calculator

import "github.com/yourusername/cost-detector/pkg/models"

// Calculator does the cost calculations
type Calculator struct {
	NodePrices map[string]float64 // Map of instance type to cost per hour
}

// NewCalculator creates a new cost calculator
func NewCalculator() *Calculator {
	return &Calculator{
		NodePrices: make(map[string]float64),
	}
}

// CalculatePodCost calculates hourly cost of a pod
func (c *Calculator) CalculatePodCost(pod *models.Pod) float64 {
	// This is simplified. In reality:
	// Cost = (Pod CPU Request / Node CPU) * Node Price Per Hour
	//      + (Pod Memory Request / Node Memory) * Node Price Per Hour
	
	// For now, rough estimate: $0.05 per CPU per hour, $0.01 per GB per hour
	cpuCost := pod.CPU * 0.05
	memoryCost := pod.Memory * 0.01
	
	return cpuCost + memoryCost
}

// CalculateHourlyCost calculates cost for multiple pods
func (c *Calculator) CalculateHourlyCost(pods []*models.Pod) float64 {
	totalCost := 0.0
	for _, pod := range pods {
		totalCost += c.CalculatePodCost(pod)
	}
	return totalCost
}
EOF
echo "  ✓ Created"

# Create alerts package (manages alerts)
echo "Creating pkg/alerts/alerts.go..."
cat > pkg/alerts/alerts.go << 'EOF'
package alerts

import "github.com/yourusername/cost-detector/pkg/models"

// Alerter handles alert logic and decisions
type Alerter struct {
	ThresholdPerHour float64 // Alert if cost exceeds this per hour
}

// NewAlerter creates a new alerter
func NewAlerter(threshold float64) *Alerter {
	return &Alerter{
		ThresholdPerHour: threshold,
	}
}

// ShouldAlert checks if we should send an alert
func (a *Alerter) ShouldAlert(costPerHour float64) bool {
	return costPerHour > a.ThresholdPerHour
}

// CreateAlert creates an alert message
func (a *Alerter) CreateAlert(team string, service string, costPerHour float64) *models.CostAlert {
	severity := "info"
	if costPerHour > a.ThresholdPerHour*2 {
		severity = "critical"
	} else if costPerHour > a.ThresholdPerHour {
		severity = "warning"
	}
	
	return &models.CostAlert{
		Team:      team,
		Service:   service,
		CostPerHr: costPerHour,
		Severity:  severity,
	}
}
EOF
echo "  ✓ Created"

# Create teams package (Teams API integration)
echo "Creating pkg/teams/teams.go..."
cat > pkg/teams/teams.go << 'EOF'
package teams

import (
	"fmt"
	"github.com/yourusername/cost-detector/pkg/models"
)

// TeamsClient sends messages to Microsoft Teams
type TeamsClient struct {
	WebhookURL string
}

// NewTeamsClient creates a new Teams client
func NewTeamsClient(webhookURL string) *TeamsClient {
	return &TeamsClient{
		WebhookURL: webhookURL,
	}
}

// SendAlert sends a cost alert to Teams
func (tc *TeamsClient) SendAlert(alert *models.CostAlert) error {
	message := fmt.Sprintf(
		"🚨 COST ALERT\nTeam: %s\nService: %s\nCost/Hr: $%.2f\nSeverity: %s",
		alert.Team,
		alert.Service,
		alert.CostPerHr,
		alert.Severity,
	)
	
	fmt.Printf("Would send to Teams: %s\n", message)
	// TODO: Actually call Teams API with webhook
	return nil
}
EOF
echo "  ✓ Created"

# Create config package (configuration management)
echo "Creating pkg/config/config.go..."
cat > pkg/config/config.go << 'EOF'
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
EOF
echo "  ✓ Created"

# Create logger package (logging)
echo "Creating pkg/logger/logger.go..."
cat > pkg/logger/logger.go << 'EOF'
package logger

import (
	"fmt"
	"time"
)

// Logger is a simple logger
type Logger struct {
	Level string // "debug", "info", "warning", "error"
}

// NewLogger creates a new logger
func NewLogger(level string) *Logger {
	return &Logger{
		Level: level,
	}
}

// Info logs info message
func (l *Logger) Info(msg string) {
	fmt.Printf("[%s] INFO: %s\n", time.Now().Format("15:04:05"), msg)
}

// Error logs error message
func (l *Logger) Error(msg string) {
	fmt.Printf("[%s] ERROR: %s\n", time.Now().Format("15:04:05"), msg)
}

// Debug logs debug message
func (l *Logger) Debug(msg string) {
	if l.Level == "debug" {
		fmt.Printf("[%s] DEBUG: %s\n", time.Now().Format("15:04:05"), msg)
	}
}
EOF
echo "  ✓ Created"

# Create main.go in cmd/cost-detector
echo "Creating cmd/cost-detector/main.go..."
cat > cmd/cost-detector/main.go << 'EOF'
package main

import (
	"fmt"
	"github.com/yourusername/cost-detector/pkg/alerts"
	"github.com/yourusername/cost-detector/pkg/calculator"
	"github.com/yourusername/cost-detector/pkg/config"
	"github.com/yourusername/cost-detector/pkg/logger"
	"github.com/yourusername/cost-detector/pkg/models"
	"github.com/yourusername/cost-detector/pkg/teams"
	"github.com/yourusername/cost-detector/pkg/watcher"
)

func main() {
	fmt.Println("🚀 Cost Detector Starting...")
	fmt.Println("")

	// Load configuration
	cfg := config.LoadConfig()
	log := logger.NewLogger(cfg.LogLevel)

	log.Info("Cost Detector initialized")
	log.Info(fmt.Sprintf("Cluster: %s", cfg.ClusterName))
	log.Info(fmt.Sprintf("Cost threshold: $%.2f/hr", cfg.CostThreshold))

	// Initialize components
	watcher := watcher.NewWatcher(cfg.ClusterName)
	calculator := calculator.NewCalculator()
	alerter := alerts.NewAlerter(cfg.CostThreshold)
	teamsClient := teams.NewTeamsClient(cfg.TeamsWebhookURL)

	// Example: Simulate some pods
	pods := []*models.Pod{
		{Name: "app-1", Namespace: "production", CPU: 2, Memory: 4},
		{Name: "app-2", Namespace: "production", CPU: 4, Memory: 8},
		{Name: "debug-app", Namespace: "debug", CPU: 16, Memory: 64}, // Expensive!
	}

	log.Info(fmt.Sprintf("Found %d pods", len(pods)))

	// Calculate total cost
	totalCost := calculator.CalculateHourlyCost(pods)
	log.Info(fmt.Sprintf("Total hourly cost: $%.2f", totalCost))

	// Check for alerts
	for _, pod := range pods {
		podCost := calculator.CalculatePodCost(pod)
		log.Info(fmt.Sprintf("Pod %s cost: $%.2f/hr", pod.Name, podCost))

		if alerter.ShouldAlert(podCost) {
			alert := alerter.CreateAlert("platform-team", pod.Name, podCost)
			teamsClient.SendAlert(alert)
		}
	}

	// Start watcher
	if err := watcher.Start(); err != nil {
		log.Error(fmt.Sprintf("Failed to start watcher: %v", err))
		return
	}

	log.Info("Cost Detector running. Press Ctrl+C to stop.")
	select {} // Run forever
}
EOF
echo "  ✓ Created"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✅ Go Module Setup Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}What we just created:${NC}"
echo "  • go.mod - Tracks your Go dependencies"
echo "  • pkg/models/ - Data structures (Pod, Alert, etc)"
echo "  • pkg/watcher/ - Watches for pod events"
echo "  • pkg/calculator/ - Calculates costs"
echo "  • pkg/alerts/ - Decides which alerts to send"
echo "  • pkg/teams/ - Sends Teams messages"
echo "  • pkg/config/ - Configuration management"
echo "  • pkg/logger/ - Simple logging"
echo "  • cmd/cost-detector/main.go - The actual app"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Update 'github.com/yourusername' to your actual GitHub username in all files"
echo "2. Test if it compiles: go build ./cmd/cost-detector"
echo "3. Run it: go run ./cmd/cost-detector/main.go"
echo ""
echo "Ready? Let's go! 🚀"
EOF
