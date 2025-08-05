package handlers

import (
	"context"

	"fulcrumproject.org/test-agent/internal/vm"
	"github.com/fulcrumproject/agent-lib-go/pkg/agent"
)

// MetricsReporter collects and returns metrics from VMs
func MetricsReporter(vmManager *vm.Manager) agent.MetricsReporter[VMServiceProperties] {
	return func(ctx context.Context) ([]agent.MetricEntry, error) {
		vms := vmManager.GetVMs()
		var metrics []agent.MetricEntry

		for _, v := range vms {
			if v.Status != vm.VMStatusSTARTED {
				// Skip VMs that aren't started
				continue
			}

			vmMetrics := []agent.MetricEntry{
				{
					TypeName:   "vm.cpu.usage",
					Value:      v.CPUUsage,
					ResourceID: v.ID,
					ExternalID: v.ID,
				},
				{
					TypeName:   "vm.memory.usage",
					Value:      v.MemoryUsage,
					ResourceID: v.ID,
					ExternalID: v.ID,
				},
				{
					TypeName:   "vm.disk.usage",
					Value:      v.DiskUsage,
					ResourceID: v.ID,
					ExternalID: v.ID,
				},
				{
					TypeName:   "vm.network.throughput",
					Value:      v.NetworkUsage,
					ResourceID: v.ID,
					ExternalID: v.ID,
				},
			}

			metrics = append(metrics, vmMetrics...)
		}

		return metrics, nil
	}
}
