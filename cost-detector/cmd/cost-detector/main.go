package main

import (
	"fmt"
	"time"
	"cost-detector/pkg/alerts"
	"cost-detector/pkg/calculator"
	"cost-detector/pkg/config"
	"cost-detector/pkg/logger"
	"cost-detector/pkg/models"
	"cost-detector/pkg/teams"
	"cost-detector/pkg/watcher"
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
	watchr := watcher.NewWatcher(cfg.ClusterName)
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
	if err := watchr.Start(); err != nil {
		log.Error(fmt.Sprintf("Failed to start watcher: %v", err))
		return
	}

	log.Info("Cost Detector running. Simulating for 5 seconds then exiting...")
	
	// Run for 5 seconds then exit (for testing)
	time.Sleep(5 * time.Second)
	
	watchr.Stop()
	log.Info("Cost Detector stopped cleanly ✅")
}