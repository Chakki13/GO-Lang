package remediators

import (
	"auto-remediation-engine/pkg/logger"
	"auto-remediation-engine/pkg/models"
)

// Remediator interface defines what a remediator can do
// WHY: Interface allows different remediation strategies
// Easy to swap implementations or add new ones
type Remediator interface {
	Remediate(pod models.Pod) (models.RemediationAction, error)
	CanHandle(pod models.Pod) bool
}

// BaseRemediator provides common functionality
type BaseRemediator struct {
	log *logger.Logger
}

// PodRestarter restarts crashed pods
type PodRestarter struct {
	BaseRemediator
}

// NewPodRestarter creates a new pod restarter
func NewPodRestarter(log *logger.Logger) *PodRestarter {
	return &PodRestarter{
		BaseRemediator{log: log},
	}
}

// Remediate restarts a pod
func (pr *PodRestarter) Remediate(pod models.Pod) (models.RemediationAction, error) {
	pr.log.Info("Restarting pod: " + pod.Name)
	
	// TODO: Actually restart the pod using Kubernetes API
	
	action := models.RemediationAction{
		Type:    "pod_restart",
		Pod:     pod,
		Status:  "completed",
		Reason:  "Pod was in CrashLoopBackOff",
		Message: "Pod restarted successfully",
	}
	
	return action, nil
}

// CanHandle checks if this remediator can handle this pod
func (pr *PodRestarter) CanHandle(pod models.Pod) bool {
	return pod.Status == "CrashLoopBackOff" || pod.Restarts > 5
}

// MemoryIncreaser increases memory requests for OOMKilled pods
type MemoryIncreaser struct {
	BaseRemediator
	incrementPercent int
}

// NewMemoryIncreaser creates a new memory increaser
func NewMemoryIncreaser(log *logger.Logger, percent int) *MemoryIncreaser {
	return &MemoryIncreaser{
		BaseRemediator:   BaseRemediator{log: log},
		incrementPercent: percent,
	}
}

// Remediate increases memory for a pod
func (mi *MemoryIncreaser) Remediate(pod models.Pod) (models.RemediationAction, error) {
	mi.log.Info("Increasing memory for pod: " + pod.Name)
	
	// TODO: Actually update the deployment and commit to Git
	
	action := models.RemediationAction{
		Type:    "memory_increase",
		Pod:     pod,
		Status:  "completed",
		Reason:  "Pod was OOMKilled",
		Message: "Memory increased by " + string(rune(mi.incrementPercent)) + "%",
	}
	
	return action, nil
}

// CanHandle checks if this remediator can handle this pod
func (mi *MemoryIncreaser) CanHandle(pod models.Pod) bool {
	return pod.Status == "OOMKilled"
}
