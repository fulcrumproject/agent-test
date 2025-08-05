package handlers

import (
	"context"
	"errors"
	"time"

	"fulcrumproject.org/test-agent/internal/vm"
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
func CreateVMHandler(vmManager *vm.Manager) agent.JobHandler[VMServiceProperties, JobResources] {
	return func(ctx context.Context, job *agent.Job[VMServiceProperties]) (*agent.JobResponse[JobResources], error) {
		if job.Service.TargetProperties == nil {
			return nil, errors.New("missing target properties")
		}
		props := *job.Service.TargetProperties

		vm, err := vmManager.CreateVM(job.Service.Name, props.CPU, props.Memory)
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
func UpdateVMHandler(vmManager *vm.Manager) agent.JobHandler[VMServiceProperties, JobResources] {
	return func(ctx context.Context, job *agent.Job[VMServiceProperties]) (*agent.JobResponse[JobResources], error) {
		if job.Service.TargetProperties == nil {
			return nil, errors.New("missing target properties")
		}
		props := *job.Service.TargetProperties
		if job.Service.ExternalID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := vmManager.UpdateVM(*job.Service.ExternalID, job.Service.Name, props.CPU, props.Memory); err != nil {
			return nil, err
		}

		return &agent.JobResponse[JobResources]{
			Resources: &JobResources{TS: time.Now()},
		}, nil
	}
}

// StartVMHandler handles VM start jobs
func StartVMHandler(vmManager *vm.Manager) agent.JobHandler[VMServiceProperties, JobResources] {
	return func(ctx context.Context, job *agent.Job[VMServiceProperties]) (*agent.JobResponse[JobResources], error) {
		if job.Service.ExternalID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := vmManager.StartVM(*job.Service.ExternalID); err != nil {
			return nil, err
		}

		return &agent.JobResponse[JobResources]{
			Resources: &JobResources{TS: time.Now()},
		}, nil
	}
}

// StopVMHandler handles VM stop jobs
func StopVMHandler(vmManager *vm.Manager) agent.JobHandler[VMServiceProperties, JobResources] {
	return func(ctx context.Context, job *agent.Job[VMServiceProperties]) (*agent.JobResponse[JobResources], error) {
		if job.Service.ExternalID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := vmManager.StopVM(*job.Service.ExternalID); err != nil {
			return nil, err
		}

		return &agent.JobResponse[JobResources]{
			Resources: &JobResources{TS: time.Now()},
		}, nil
	}
}

// DeleteVMHandler handles VM delete jobs
func DeleteVMHandler(vmManager *vm.Manager) agent.JobHandler[VMServiceProperties, JobResources] {
	return func(ctx context.Context, job *agent.Job[VMServiceProperties]) (*agent.JobResponse[JobResources], error) {
		if job.Service.ExternalID == nil {
			return nil, errors.New("missing externalId")
		}

		if err := vmManager.DeleteVM(*job.Service.ExternalID); err != nil {
			return nil, err
		}

		return &agent.JobResponse[JobResources]{
			Resources: &JobResources{TS: time.Now()},
		}, nil
	}
}
