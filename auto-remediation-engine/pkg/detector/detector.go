package detector

import (
	"auto-remediation-engine/pkg/logger"
	"auto-remediation-engine/pkg/models"
)

// Detector watches the cluster for issues
// WHY: We need to continuously monitor for problems
// This is the "eyes" of the remediation engine
type Detector struct {
	log *logger.Logger
}

// NewDetector creates a new detector
func NewDetector(log *logger.Logger) *Detector {
	return &Detector{
		log: log,
	}
}

// DetectPodCrashes checks for pods in CrashLoopBackOff
// WHY: CrashLoopBackOff indicates the pod keeps crashing
// We can restart it to break the loop
func (d *Detector) DetectPodCrashes(pods []models.Pod) []models.Pod {
	d.log.Info("Scanning for crashed pods...")
	var crashed []models.Pod
	
	for _, pod := range pods {
		if pod.Status == "CrashLoopBackOff" || pod.Restarts > 5 {
			d.log.Warning("Found crashed pod: " + pod.Name)
			crashed = append(crashed, pod)
		}
	}
	
	return crashed
}

// DetectOOMKilled checks for out-of-memory killed pods
// WHY: OOMKilled means the pod needs more memory
// We can increase memory and prevent future crashes
func (d *Detector) DetectOOMKilled(pods []models.Pod) []models.Pod {
	d.log.Info("Scanning for OOMKilled pods...")
	var oomPods []models.Pod
	
	for _, pod := range pods {
		if pod.Status == "OOMKilled" {
			d.log.Warning("Found OOMKilled pod: " + pod.Name)
			oomPods = append(oomPods, pod)
		}
	}
	
	return oomPods
}

// DetectSecurityIssues checks for security problems
// WHY: Security vulnerabilities can cause data breaches
// We should fix these automatically
func (d *Detector) DetectSecurityIssues(pods []models.Pod) []models.SecurityIssue {
	d.log.Info("Scanning for security issues...")
	var issues []models.SecurityIssue
	
	// TODO: Implement actual security context checking
	// For now, this is a placeholder
	
	return issues
}

// DetectMissingProbes checks for missing health checks
// WHY: Missing probes mean Kubernetes can't detect dead pods
// We should alert but NOT auto-inject (needs manual review)
func (d *Detector) DetectMissingProbes(pods []models.Pod) []models.ProbeIssue {
	d.log.Info("Scanning for missing probes...")
	var issues []models.ProbeIssue
	
	// TODO: Implement actual probe checking
	
	return issues
}

// DetectExpiringCerts checks certificate expiry
// WHY: Expired certs cause outages
// Early warning allows teams to prepare
func (d *Detector) DetectExpiringCerts(threshold int) []models.CertificateAlert {
	d.log.Info("Scanning for expiring certificates...")
	var alerts []models.CertificateAlert
	
	// TODO: Implement actual certificate checking
	
	return alerts
}
