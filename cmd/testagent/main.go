package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fulcrumproject/agent-lib-go/pkg/fulcrumcli"
	"github.com/fulcrumproject/agent-lib-go/pkg/stdagent"
	"github.com/fulcrumproject/test-agent/internal/config"
	"github.com/fulcrumproject/test-agent/internal/vm"
	"github.com/fulcrumproject/test-agent/internal/vmagent"
	"github.com/fulcrumproject/utils/confbuilder"
)

type AgentRemoteConfig struct{}

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := confbuilder.New(config.DefaultConfig()).
		EnvPrefix(config.EnvPrefix).
		EnvFiles(".env").
		File(configPath).
		Build()
	if err != nil {
		slog.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}

	log.Printf("Starting test agent with Fulcrum API at %s", cfg.FulcrumAPIURL)

	client := fulcrumcli.NewHTTPClient(cfg.FulcrumAPIURL, cfg.AgentToken)
	vmManager := vm.NewManager(cfg.ErrorRate, cfg.OperationDelayMin, cfg.OperationDelayMax)

	// Create stdagent with configuration options
	testAgent := vmagent.NewVMAgent(vmManager, client,
		stdagent.WithHealthInterval(60*time.Second),
		stdagent.WithJobPollInterval(cfg.JobPollInterval),
		stdagent.WithMetricsReportInterval(cfg.MetricReportInterval),
	)

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start the agent
	if err := testAgent.Run(ctx); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	log.Printf("Test agent started successfully (Agent ID: %s)", testAgent.GetAgentID())
	log.Printf("Press Ctrl+C to stop the agent")

	// Start periodic status reporting
	go startStatusReporter(ctx, vmManager, testAgent)

	// Wait for termination signal
	<-sigCh
	log.Println("Received shutdown signal")

	// Create a context with timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shut down the agent
	if err := testAgent.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}

	// Display final job statistics
	processed, succeeded, failed, unsupported := testAgent.GetJobStats()
	log.Printf("Final Job Statistics: Processed: %d, Succeeded: %d, Failed: %d, Unsupported: %d", processed, succeeded, failed, unsupported)

	// Display final VM status counts
	statusCounts := vmManager.GetStatusCounts()
	if len(statusCounts) > 0 {
		log.Printf("Final VM Statuss: %v", statusCounts)
	}
	log.Printf("Agent uptime: %s", testAgent.GetUptime().Round(time.Second))
}

// startStatusReporter starts a background goroutine to periodically display VM status and job statistics
func startStatusReporter(ctx context.Context, vmManager *vm.Manager, testAgent *vmagent.VMAgent) {
	displayTicker := time.NewTicker(30 * time.Second)
	defer displayTicker.Stop()
	for {
		select {
		case <-displayTicker.C:
			// Display VM status counts
			statusCounts := vmManager.GetStatusCounts()
			log.Printf("VM Statuss: %v", statusCounts)

			// Display job statistics
			processed, succeeded, failed, unsupported := testAgent.GetJobStats()
			log.Printf("Jobs: Processed: %d, Succeeded: %d, Failed: %d, Unsupported: %d",
				processed, succeeded, failed, unsupported)
		case <-ctx.Done():
			return
		}
	}
}
