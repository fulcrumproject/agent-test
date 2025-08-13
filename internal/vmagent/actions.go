package vmagent

import (
	"context"
	"errors"
	"time"

	"github.com/fulcrumproject/agent-lib-go/pkg/agent"
)

// VMServiceProperties represents the properties of a VM service
type VMServiceProperties struct {
	CPU    int `json:"cpu"`
	Memory int `json:"memory"`
}

// JobResources represents the resources in a job response
type JobResources struct {
	TS time.Time `json:"ts"`
}

// CreateVMHandler handles VM creation jobs
func (v *VMAgent) CreateVMHandler() agent.JobHandler[VMServiceProperties, JobResources] {
	return func(ctx context.Context, job *agent.Job[VMServiceProperties]) (*agent.JobResponse[JobResources], error) {
		if job.Params == nil {
			return nil, errors.New("missing target properties")
		}

		vm, err := v.vmManager.CreateVM(job.Service.Name, job.Params.CPU, job.Params.Memory)
		if err != nil {
			return nil, err
		}

		return &agent.JobResponse[JobResources]{
			Resources:  &JobResources{TS: time.Now()},
			ExternalID: &vm.ID,
		}, nil
	}
}

// UpdateVMHandler handles VM update jobs
func (v *VMAgent) UpdateVMHandler() agent.JobHandler[VMServiceProperties, JobResources] {
	return func(ctx context.Context, job *agent.Job[VMServiceProperties]) (*agent.JobResponse[JobResources], error) {
		if job.Params == nil {
			return nil, errors.New("missing target properties")
		}
		if job.Service.ExternalID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := v.vmManager.UpdateVM(*job.Service.ExternalID, job.Service.Name, job.Params.CPU, job.Params.Memory); err != nil {
			return nil, err
		}

		return &agent.JobResponse[JobResources]{
			Resources: &JobResources{TS: time.Now()},
		}, nil
	}
}

// StartVMHandler handles VM start jobs
func (v *VMAgent) StartVMHandler() agent.JobHandler[VMServiceProperties, JobResources] {
	return func(ctx context.Context, job *agent.Job[VMServiceProperties]) (*agent.JobResponse[JobResources], error) {
		if job.Service.ExternalID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := v.vmManager.StartVM(*job.Service.ExternalID); err != nil {
			return nil, err
		}

		return &agent.JobResponse[JobResources]{
			Resources: &JobResources{TS: time.Now()},
		}, nil
	}
}

// StopVMHandler handles VM stop jobs
func (v *VMAgent) StopVMHandler() agent.JobHandler[VMServiceProperties, JobResources] {
	return func(ctx context.Context, job *agent.Job[VMServiceProperties]) (*agent.JobResponse[JobResources], error) {
		if job.Service.ExternalID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := v.vmManager.StopVM(*job.Service.ExternalID); err != nil {
			return nil, err
		}

		return &agent.JobResponse[JobResources]{
			Resources: &JobResources{TS: time.Now()},
		}, nil
	}
}

// DeleteVMHandler handles VM delete jobs
func (v *VMAgent) DeleteVMHandler() agent.JobHandler[VMServiceProperties, JobResources] {
	return func(ctx context.Context, job *agent.Job[VMServiceProperties]) (*agent.JobResponse[JobResources], error) {
		if job.Service.ExternalID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := v.vmManager.DeleteVM(*job.Service.ExternalID); err != nil {
			return nil, err
		}

		return &agent.JobResponse[JobResources]{
			Resources: &JobResources{TS: time.Now()},
		}, nil
	}
}
