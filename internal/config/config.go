package config

import (
	"time"
)

const (
	EnvPrefix = "TESTAGENT_"
)

// Config holds the configuration for the test agent
type Config struct {
	// Agent authentication
	AgentToken string `json:"agentToken" validate:"required" env:"AGENT_TOKEN"` // Authentication token for the agent

	// Fulcrum Core API connection
	FulcrumAPIURL string `json:"fulcrumApiUrl" validate:"required" env:"FULCRUM_API_URL"`

	// Simulation parameters
	VMUpdateInterval     time.Duration `json:"vmUpdateInterval" env:"VM_OPERATION_INTERVAL"`      // How often to perform VM operations
	JobPollInterval      time.Duration `json:"jobPollInterval" env:"JOB_POLL_INTERVAL"`           // How often to poll for jobs
	MetricReportInterval time.Duration `json:"metricReportInterval" env:"METRIC_REPORT_INTERVAL"` // How often to report metrics

	// Simulation behavior
	OperationDelayMin time.Duration `json:"operationDelayMin" env:"OPERATION_DELAY_MIN"` // Minimum time for operation
	OperationDelayMax time.Duration `json:"operationDelayMax" env:"OPERATION_DELAY_MAX"` // Maximum time for operation
	ErrorRate         float64       `json:"errorRate" env:"ERROR_RATE"`                  // Probability of operation failure (0.0-1.0)
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		AgentToken:           "", // Must be provided
		FulcrumAPIURL:        "http://localhost:3000",
		VMUpdateInterval:     5 * time.Second,
		JobPollInterval:      5 * time.Second,
		MetricReportInterval: 30 * time.Second,
		OperationDelayMin:    2 * time.Second,
		OperationDelayMax:    10 * time.Second,
		ErrorRate:            0.05, // 5% chance of failure
	}
}
