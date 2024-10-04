package status

import (
	// "encoding/json"
	// "log"
	// "os"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// const statusFilePath = "/opt/k8s-agent/status.json"
const agentEnvPath = "/etc/default/k8s-agent"

var mu sync.RWMutex

// List of states
const (
	Ready      = "Ready"
	Preparing  = "Preparing"
	Deploying  = "Deploying"
	Running    = "Running"
	Restarting = "Restarting"
	Error      = "Error"
)

// List of node types
const (
	Bootstrap = "bootstrap"
	Server    = "server"
	Agent     = "agent"
)

// Status struct holds the state of the agent and related information
type Status struct {
	// mu sync.RWMutex `json:"-"` // Mutex is ignored during JSON marshalling

	State    string `json:"state"`       // Current agent state
	Role     string `json:"node_role"`   // Node role (bootstrap/server/agent)
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

// Initialize state
func init() {
	mu.Lock()
	defer mu.Unlock()

	// Load env file
	godotenv.Load(agentEnvPath)

	// Set from env
	agentStatus = Status{
		State:    Ready,
		Role:     os.Getenv("NODE_ROLE"),
		Version:  os.Getenv("AGENT_VERSION"),
		Release:  os.Getenv("RKE2_RELEASE"),
		PublicIP: os.Getenv("PUBLIC_IP"),
		ApiIP:    os.Getenv("VIRTUAL_IP") + "/" + os.Getenv("VIRTUAL_CIDR"),
		Region:   "",
		GpuCount: 0,
	}

	// 	// Try to load from file, ignoring errors
	// 	data, _ := os.ReadFile(statusFilePath)
	// 	json.Unmarshal(data, &agentStatus)

	// 	// Write status to file in case it's missing/corrupted
	// 	data, err := json.Marshal(agentStatus)
	// 	if err != nil {
	// 		log.Fatal("Internal error marshaling status")
	// 	}

	// err = os.WriteFile(statusFilePath, data, 0600)
	//
	//	if err != nil {
	//		log.Fatal("Fatal error saving status file")
	//	}
}
