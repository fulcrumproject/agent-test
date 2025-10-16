package vmagent

import (
	"context"
	"fmt"
	"log"

	"github.com/fulcrumproject/agent-lib-go/pkg/agent"
)

// MetricsReporter collects and returns metrics from VMs
func (v *VMAgent) MetricsReporter() agent.MetricsReporter[VMProperties, VMResources] {
	return func(ctx context.Context, service *agent.Service[VMProperties, VMResources]) ([]agent.MetricEntry, error) {
		var metrics []agent.MetricEntry

		if service == nil || service.AgentInstanceID == nil {
			return nil, fmt.Errorf("service is nil or externalID is nil")
		}

		vmInstance, ok := v.vmManager.GetVM(*service.AgentInstanceID)
		if !ok {
			log.Printf("ignoring vm for metrics: instance not found for externalID: %s", *service.AgentInstanceID)
			return metrics, nil
		}

		metrics = append(metrics,
			agent.MetricEntry{
				TypeName:        "vm.cpu.usage",
				Value:           vmInstance.CPUUsage,
				ResourceID:      vmInstance.ID,
				AgentInstanceID: vmInstance.ID,
			},
			agent.MetricEntry{
				TypeName:        "vm.memory.usage",
				Value:           vmInstance.MemoryUsage,
				ResourceID:      vmInstance.ID,
				AgentInstanceID: vmInstance.ID,
			},
			agent.MetricEntry{
				TypeName:        "vm.disk.usage",
				Value:           vmInstance.DiskUsage,
				ResourceID:      vmInstance.ID,
				AgentInstanceID: vmInstance.ID,
			},
			agent.MetricEntry{
				TypeName:        "vm.network.throughput",
				Value:           vmInstance.NetworkUsage,
				ResourceID:      vmInstance.ID,
				AgentInstanceID: vmInstance.ID,
			},
		)

		return metrics, nil
	}
}
