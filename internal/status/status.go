package status

import (
	// "encoding/json"
	// "log"
	// "os"
	"sync"
	// "github.com/joho/godotenv"
)

// const statusFilePath = "/opt/k8s-agent/status.json"

var mu sync.RWMutex

// List of states
const (
	Preparing  = "Preparing"
	Ready      = "Ready"
	Deploying  = "Deploying"
	Running    = "Running"
	Restarting = "Restarting"
	Error      = "Error"
)

// List of node types
const (
	Server = "server"
	Agent  = "agent"
)

// Status struct holds the state of the agent and related information
type Status struct {
	// mu sync.RWMutex `json:"-"` // Mutex is ignored during JSON marshalling

	State    string `json:"state"`       // Current agent state
	Type     string `json:"node_type"`   // Node type (server/agent)
	Version  string `json:"version"`     // Agent version
	Release  string `json:"k8s_release"` // K8s release
	PublicIP string `json:"public_ip"`   // Node public IP
	ApiIP    string `json:"api_ip"`      // Control plane IP/VIP
	Region   string `json:"region"`      // Node region
	GpuCount int    `json:"gpu_count"`   // Node GPU count
}

var agentStatus = Status{}

// GetStatus returns a copy of the current status with a read lock
func GetStatus() Status {
	mu.RLock()
	defer mu.RUnlock()
	return agentStatus
}

// UpdateStatus safely updates the status with a write lock and persists it
func UpdateStatus(newStatus Status) error {
	mu.Lock()
	defer mu.Unlock()
	agentStatus = newStatus

	// data, err := json.Marshal(newStatus)
	// if err != nil {
	// 	return err
	// }

	// return os.WriteFile(statusFilePath, data, 0600)
	return nil
}

// // Initialize state
// func init() {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	// Set defaults
// 	agentStatus = Status{
// 		State: Preparing,
// 		Type:  Server,
// 	}

// 	// Try to load from file, ignoring errors
// 	data, _ := os.ReadFile(statusFilePath)
// 	json.Unmarshal(data, &agentStatus)

// 	// Write status to file in case it's missing/corrupted
// 	data, err := json.Marshal(agentStatus)
// 	if err != nil {
// 		log.Fatal("Internal error marshaling status")
// 	}

// 	err = os.WriteFile(statusFilePath, data, 0600)
// 	if err != nil {
// 		log.Fatal("Fatal error saving status file")
// 	}
// }
