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
