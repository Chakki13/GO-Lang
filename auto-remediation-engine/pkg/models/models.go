package models

import "time"

// Pod represents a Kubernetes pod we're monitoring
// WHY: We need to track pod state, resources, and health
type Pod struct {
	Name       string    // Pod name in cluster
	Namespace  string    // Which namespace it's in
	NodeName   string    // Which node it's running on
	Status     string    // Running, Failed, CrashLoopBackOff, etc
	Restarts   int       // How many times it restarted
	CPU        string    // CPU request (e.g., "500m")
	Memory     string    // Memory request (e.g., "512Mi")
	CreatedAt  time.Time // When pod was created
}

// RemediationAction represents an action the system can take
// WHY: We need to know what action was taken, when, and why
type RemediationAction struct {
	ID            string    // Unique ID for this action
	Type          string    // "pod_restart", "memory_increase", etc
	Pod           Pod       // Which pod this affected
	Status        string    // "pending", "executing", "completed", "failed"
	Reason        string    // Why we took this action
	Message       string    // Human-readable description
	GitCommit     string    // Git commit hash if we modified something
	ExecutedAt    time.Time // When we executed it
	RequiresApproval bool   // Did this need manual approval?
	ApprovedBy    string    // Who approved it (if applicable)
}

// CertificateAlert represents an expiring certificate
// WHY: Certificates are critical; we need to track expiry and alert early
type CertificateAlert struct {
	Name          string    // Certificate name
	ExpiresIn     int       // Days until expiry
	ExpiresAt     time.Time // Exact expiry date
	Service       string    // Which service uses this cert
	Namespace     string    // Which namespace
	Severity      string    // "info", "warning", "critical"
}

// SecurityIssue represents a security vulnerability found
// WHY: Security is non-negotiable; we track and fix security issues
type SecurityIssue struct {
	Pod          Pod       // Which pod has the issue
	IssueType    string    // "running_as_root", "missing_security_context", etc
	Description  string    // What's the problem?
	Remediation  string    // How to fix it
	CanAutoFix   bool      // Can we fix this automatically?
}

// ProbeIssue represents a missing health check probe
// WHY: Missing probes = Kubernetes can't detect dead pods
type ProbeIssue struct {
	Pod        Pod    // Which pod is affected
	ProbeType  string // "liveness" or "readiness"
	Message    string // Recommendation for the team
}

// ClusterMetrics represents overall cluster health
// WHY: We need to understand if the cluster is under stress
type ClusterMetrics struct {
	TotalPods      int
	CrashedPods    int
	OOMKilledPods  int
	RestartingPods int
	TimestampAt    time.Time
}
