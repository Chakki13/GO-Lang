package teams

import (
	"fmt"
	"cost-detector/pkg/models"
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
