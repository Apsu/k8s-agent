package status

import (
	// "encoding/json"
	// "log"
	"fmt"
	// "os"
	"sync"

	// "github.com/joho/godotenv"
	"github.com/caarlos0/env/v11"
)

// const statusFilePath = "/opt/k8s-agent/status.json"
const agentEnvPath = "/etc/default/k8s-agent"

var mu sync.RWMutex

// List of states
const (
	Initialized = "Initialized"
	Preparing   = "Preparing"
	Deploying   = "Deploying"
	Restarting  = "Restarting"
	Ready       = "Ready"
	Error       = "Error"
)

// List of node types
const (
	Bootstrap = "bootstrap"
	Server    = "server"
	Agent     = "agent"
)

// Status struct holds the state of the agent and related information
// mu sync.RWMutex `json:"-"` // Mutex is ignored during JSON marshalling
type Status struct {
	State     string `env:"STATE" envdefault:"Initialized"` // Current agent state
	PublicIP  string `env:"PUBLIC_IP"`                      // Node public IP
	ApiIP     string `env:"API_IP"`                         // Control plane IP/VIP
	ApiCIDR   string `env:"API_CIDR"`                       // Control plane CIDR
	Interface string `env:"INTERFACE"`                      // Node interface
	Release   string `env:"RKE2_RELEASE"`                   // K8s release
	Region    string `env:"NODE_REGION"`                    // Node region
	Role      string `env:"NODE_ROLE"`                      // Node role (bootstrap/server/agent)
	GpuCount  int    `env:"GPU_COUNT"`                      // Node GPU count
	Version   string `env:"AGENT_VERSION"`                  // Agent version
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

	// file, err := os.Create(agentEnvPath)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// // Write each field as "KEY=value"
	// fmt.Fprintf(file, "STATE=%s\n", agentStatus.State)
	// fmt.Fprintf(file, "PUBLIC_IP=%s\n", agentStatus.PublicIP)
	// fmt.Fprintf(file, "API_IP=%s\n", agentStatus.ApiIP)
	// fmt.Fprintf(file, "API_CIDR=%s\n", agentStatus.ApiCIDR)
	// fmt.Fprintf(file, "INTERFACE=%s\n", agentStatus.Interface)
	// fmt.Fprintf(file, "RKE2_RELEASE=%s\n", agentStatus.Release)
	// fmt.Fprintf(file, "NODE_REGION=%s\n", agentStatus.Region)
	// fmt.Fprintf(file, "NODE_ROLE=%s\n", agentStatus.Role)
	// fmt.Fprintf(file, "GPU_COUNT=%d\n", agentStatus.GpuCount)
	// fmt.Fprintf(file, "AGENT_VERSION=%s\n", agentStatus.Version)

	return nil
}

// Initialize state
func init() {
	mu.Lock()
	defer mu.Unlock()

	// Load env
	env.Parse(agentStatus)
	fmt.Print(agentStatus)
}
