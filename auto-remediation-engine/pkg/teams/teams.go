package teams

import (
	"fmt"
	"auto-remediation-engine/pkg/logger"
	"auto-remediation-engine/pkg/models"
)

// TeamsClient sends alerts to Microsoft Teams
// WHY: Teams is where your team communicates
// Real-time alerts ensure issues are seen immediately
type TeamsClient struct {
	webhookURL string
	log        *logger.Logger
}

// NewTeamsClient creates a new Teams client
func NewTeamsClient(webhookURL string, log *logger.Logger) *TeamsClient {
	return &TeamsClient{
		webhookURL: webhookURL,
		log:        log,
	}
}

// SendRemediationAlert sends an alert about remediation action
func (tc *TeamsClient) SendRemediationAlert(action models.RemediationAction) error {
	message := fmt.Sprintf(
		"✅ AUTO-REMEDIATED\nType: %s\nPod: %s\nNamespace: %s\nReason: %s",
		action.Type,
		action.Pod.Name,
		action.Pod.Namespace,
		action.Reason,
	)
	
	tc.log.Info("Would send to Teams: " + message)
	// TODO: Actually send to Teams webhook
	
	return nil
}

// SendCertAlert sends an alert about expiring certificate
func (tc *TeamsClient) SendCertAlert(alert models.CertificateAlert) error {
	message := fmt.Sprintf(
		"⚠️  CERT EXPIRING\nService: %s\nExpires in: %d days",
		alert.Service,
		alert.ExpiresIn,
	)
	
	tc.log.Info("Would send to Teams: " + message)
	// TODO: Actually send to Teams webhook
	
	return nil
}

// SendSecurityAlert sends an alert about security issues
func (tc *TeamsClient) SendSecurityAlert(issue models.SecurityIssue) error {
	message := fmt.Sprintf(
		"🔒 SECURITY ISSUE FIXED\nPod: %s\nIssue: %s",
		issue.Pod.Name,
		issue.IssueType,
	)
	
	tc.log.Info("Would send to Teams: " + message)
	// TODO: Actually send to Teams webhook
	
	return nil
}
