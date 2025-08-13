package vmagent

import (
	"context"

	"github.com/fulcrumproject/agent-lib-go/pkg/agent"
)

// healthHandler handles health events and updates VM resources
func (v *VMAgent) HealthHandler() agent.HealthHandler {
	return func(ctx context.Context) error {
		v.vmManager.UpdateVMResources()
		return nil
	}
}
