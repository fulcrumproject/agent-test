package vmagent

import (
	"context"
	"log"
	"time"

	"github.com/fulcrumproject/agent-lib-go/pkg/agent"
	"github.com/fulcrumproject/agent-lib-go/pkg/fulcrumcli"
	"github.com/fulcrumproject/agent-lib-go/pkg/stdagent"
	"github.com/fulcrumproject/test-agent/internal/vm"
)

type RemoteConfiguration struct {
	Environment string `json:"environment"`
}

type VMAgent struct {
	stdAgent  *stdagent.Agent
	vmManager *vm.Manager
}

func NewVMAgent(vmManager *vm.Manager, client *fulcrumcli.HTTPClient, opts ...stdagent.AgentOption) *VMAgent {
	stdAgent, err := stdagent.New(client, opts...)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}
	vmAgent := &VMAgent{
		stdAgent:  stdAgent,
		vmManager: vmManager,
	}

	// Register connection handler
	stdAgent.OnConnect(agent.ConnectHandlerWrapper(vmAgent.ConnectHandler))

	// Register health handler (with VM resource updating)
	stdAgent.OnHealth(vmAgent.HealthHandler())

	// Register job handlers
	stdAgent.OnJob(agent.JobActionServiceCreate, agent.JobHandlerWrapper(vmAgent.CreateVMHandler()))
	stdAgent.OnJob(agent.JobActionServiceUpdate, agent.JobHandlerWrapper(vmAgent.UpdateVMHandler()))
	stdAgent.OnJob(agent.JobActionServiceStart, agent.JobHandlerWrapper(vmAgent.StartVMHandler()))
	stdAgent.OnJob(agent.JobActionServiceStop, agent.JobHandlerWrapper(vmAgent.StopVMHandler()))
	stdAgent.OnJob(agent.JobActionServiceDelete, agent.JobHandlerWrapper(vmAgent.DeleteVMHandler()))

	// Register metrics reporter
	stdAgent.OnMetrics(vmAgent.MetricsReporter())

	return vmAgent
}

func (v *VMAgent) ConnectHandler(ctx context.Context, info *agent.AgentInfo[RemoteConfiguration]) error {
	return nil
}

func (v *VMAgent) GetAgentID() string {
	return v.stdAgent.GetAgentID()
}

func (v *VMAgent) Run(ctx context.Context) error {
	return v.stdAgent.Run(ctx)
}

func (v *VMAgent) Shutdown(ctx context.Context) error {
	return v.stdAgent.Shutdown(ctx)
}

func (v *VMAgent) GetUptime() time.Duration {
	return v.stdAgent.GetUptime()
}

func (v *VMAgent) GetJobStats() (int, int, int, int) {
	return v.stdAgent.GetJobStats()
}
