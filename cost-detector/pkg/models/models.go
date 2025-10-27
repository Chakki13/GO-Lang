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
