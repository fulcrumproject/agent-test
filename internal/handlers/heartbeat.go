package handlers

import (
	"context"

	"fulcrumproject.org/test-agent/internal/vm"
	"github.com/fulcrumproject/agent-lib-go/pkg/agent"
)

// heartbeatHandler handles heartbeat events and updates VM resources
func HeartbeatHandler(vmManager *vm.Manager) agent.HeartbeatHandler {
	return func(ctx context.Context) error {
		vmManager.UpdateVMResources()
		return nil
	}
}
