package alerts

import "cost-detector/pkg/models"

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
