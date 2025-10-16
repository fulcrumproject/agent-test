package vmagent

import (
	"context"
	"errors"
	"time"

	"github.com/fulcrumproject/agent-lib-go/pkg/agent"
)

// VMProperties represents the properties of a VM service
type VMProperties struct {
	CPU    int `json:"cpu"`
	Memory int `json:"memory"`
}

// VMResources represents the resources in a job response
type VMResources struct {
	TS time.Time `json:"ts"`
}

// CreateVMHandler handles VM creation jobs
func (v *VMAgent) CreateVMHandler() agent.JobHandler[VMProperties, VMProperties, VMResources] {
	return func(ctx context.Context, job *agent.Job[VMProperties, VMProperties, VMResources]) (*agent.JobResponse[VMResources], error) {
		if job.Params == nil {
			return nil, errors.New("missing target properties")
		}

		vm, err := v.vmManager.CreateVM(job.Service.Name, job.Params.CPU, job.Params.Memory)
		if err != nil {
			return nil, err
		}

		return &agent.JobResponse[VMResources]{
			AgentData:       &VMResources{TS: time.Now()},
			AgentInstanceID: &vm.ID,
		}, nil
	}
}

// UpdateVMHandler handles VM update jobs
func (v *VMAgent) UpdateVMHandler() agent.JobHandler[VMProperties, VMProperties, VMResources] {
	return func(ctx context.Context, job *agent.Job[VMProperties, VMProperties, VMResources]) (*agent.JobResponse[VMResources], error) {
		if job.Params == nil {
			return nil, errors.New("missing target properties")
		}
		if job.Service.AgentInstanceID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := v.vmManager.UpdateVM(*job.Service.AgentInstanceID, job.Service.Name, job.Params.CPU, job.Params.Memory); err != nil {
			return nil, err
		}

		return &agent.JobResponse[VMResources]{
			AgentData: &VMResources{TS: time.Now()},
		}, nil
	}
}

// StartVMHandler handles VM start jobs
func (v *VMAgent) StartVMHandler() agent.JobHandler[VMProperties, VMProperties, VMResources] {
	return func(ctx context.Context, job *agent.Job[VMProperties, VMProperties, VMResources]) (*agent.JobResponse[VMResources], error) {
		if job.Service.AgentInstanceID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := v.vmManager.StartVM(*job.Service.AgentInstanceID); err != nil {
			return nil, err
		}

		return &agent.JobResponse[VMResources]{
			AgentData: &VMResources{TS: time.Now()},
		}, nil
	}
}

// StopVMHandler handles VM stop jobs
func (v *VMAgent) StopVMHandler() agent.JobHandler[VMProperties, VMProperties, VMResources] {
	return func(ctx context.Context, job *agent.Job[VMProperties, VMProperties, VMResources]) (*agent.JobResponse[VMResources], error) {
		if job.Service.AgentInstanceID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := v.vmManager.StopVM(*job.Service.AgentInstanceID); err != nil {
			return nil, err
		}

		return &agent.JobResponse[VMResources]{
			AgentData: &VMResources{TS: time.Now()},
		}, nil
	}
}

// DeleteVMHandler handles VM delete jobs
func (v *VMAgent) DeleteVMHandler() agent.JobHandler[VMProperties, VMProperties, VMResources] {
	return func(ctx context.Context, job *agent.Job[VMProperties, VMProperties, VMResources]) (*agent.JobResponse[VMResources], error) {
		if job.Service.AgentInstanceID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := v.vmManager.DeleteVM(*job.Service.AgentInstanceID); err != nil {
			return nil, err
		}

		return &agent.JobResponse[VMResources]{
			AgentData: &VMResources{TS: time.Now()},
		}, nil
	}
}
