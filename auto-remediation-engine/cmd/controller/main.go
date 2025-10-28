package main

import (
	"fmt"
	"auto-remediation-engine/pkg/config"
	"auto-remediation-engine/pkg/detector"
	"auto-remediation-engine/pkg/logger"
	"auto-remediation-engine/pkg/models"
	"auto-remediation-engine/pkg/remediators"
	"auto-remediation-engine/pkg/teams"
)

func main() {
	fmt.Println("🚀 Auto-Remediation Engine Starting...")
	fmt.Println("")

	// Load configuration from environment variables
	cfg := config.LoadConfig()
	log := logger.NewLogger(cfg.LogLevel)

	log.Info("Configuration loaded")
	log.Info(fmt.Sprintf("Cluster: %s", cfg.ClusterName))
	log.Info(fmt.Sprintf("Region: %s", cfg.ClusterRegion))

	// Initialize components
	det := detector.NewDetector(log)
	restarter := remediators.NewPodRestarter(log)
	memoryIncreaser := remediators.NewMemoryIncreaser(log, cfg.MemoryIncrementPercent)
	teamsClient := teams.NewTeamsClient(cfg.TeamsWebhookURL, log)

	// Simulate some pods for testing
	testPods := []models.Pod{
		{
			Name:      "web-app-1",
			Namespace: "production",
			Status:    "Running",
			Restarts:  0,
			CPU:       "500m",
			Memory:    "512Mi",
		},
		{
			Name:      "crashed-app",
			Namespace: "production",
			Status:    "CrashLoopBackOff",
			Restarts:  7,
			CPU:       "500m",
			Memory:    "512Mi",
		},
		{
			Name:      "memory-hog",
			Namespace: "production",
			Status:    "OOMKilled",
			Restarts:  3,
			CPU:       "1000m",
			Memory:    "1024Mi",
		},
	}

	log.Info(fmt.Sprintf("Scanning %d pods for issues...", len(testPods)))

	// Detect crashed pods
	crashedPods := det.DetectPodCrashes(testPods)
	if len(crashedPods) > 0 {
		log.Info(fmt.Sprintf("Found %d crashed pods, remediating...", len(crashedPods)))
		for _, pod := range crashedPods {
			if restarter.CanHandle(pod) {
				action, err := restarter.Remediate(pod)
				if err == nil {
					teamsClient.SendRemediationAlert(action)
				}
			}
		}
	}

	// Detect OOMKilled pods
	oomPods := det.DetectOOMKilled(testPods)
	if len(oomPods) > 0 {
		log.Info(fmt.Sprintf("Found %d OOMKilled pods, remediating...", len(oomPods)))
		for _, pod := range oomPods {
			if memoryIncreaser.CanHandle(pod) {
				action, err := memoryIncreaser.Remediate(pod)
				if err == nil {
					teamsClient.SendRemediationAlert(action)
				}
			}
		}
	}

	log.Info("Auto-Remediation Engine ready. Waiting for pod events...")
	log.Info("(In production, this would watch Kubernetes API)")
}
