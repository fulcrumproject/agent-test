package vm

import "time"

// VMStatus represents the possible statuss of a VM
type VMStatus string

const (
	VMStatusCREATED VMStatus = "CREATED"
	VMStatusSTARTED VMStatus = "STARTED"
	VMStatusSTOPPED VMStatus = "STOPPED"
	VMStatusDELETED VMStatus = "DELETED"
)

// VM represents a simulated virtual machine
type VM struct {
	ID           string
	Name         string
	Status       VMStatus
	CreatedAt    time.Time
	CPU          int
	Memory       int
	CPUUsage     float64 // Simulated CPU usage (0-100%)
	MemoryUsage  float64 // Simulated memory usage (0-100%)
	DiskUsage    float64 // Simulated disk usage (0-100%)
	NetworkUsage float64 // Simulated network throughput (Mbps)
	ErrorMessage string  // Contains error message if Status is ERROR
}
