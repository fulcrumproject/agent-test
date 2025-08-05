package vm

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Manager handles the simulation of VM lifecycles
type Manager struct {
	vms           map[string]*VM
	mutex         sync.RWMutex
	errorRate     float64
	delayRangeMin time.Duration
	delayRangeMax time.Duration
}

// NewManager creates a new VM manager
func NewManager(errorRate float64, delayRangeMin, delayRangeMax time.Duration) *Manager {
	return &Manager{
		vms:           make(map[string]*VM),
		errorRate:     errorRate,
		delayRangeMin: delayRangeMin,
		delayRangeMax: delayRangeMax,
	}
}

// GetVMs returns all managed VMs
func (m *Manager) GetVMs() []*VM {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	vms := make([]*VM, 0, len(m.vms))
	for _, vm := range m.vms {
		vms = append(vms, vm)
	}
	return vms
}

// GetVM returns a VM by ID
func (m *Manager) GetVM(id string) (*VM, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	vm, exists := m.vms[id]
	return vm, exists
}

// CreateVM starts the VM creation process
func (m *Manager) CreateVM(name string, cpu int, memory int) (*VM, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	id := fmt.Sprintf("vm-%s", uuid.New())

	vm := &VM{
		ID:           id,
		Name:         name,
		Status:       VMStatusCREATED,
		CreatedAt:    time.Now(),
		CPUUsage:     0,
		MemoryUsage:  0,
		DiskUsage:    0,
		NetworkUsage: 0,
	}

	delay := m.randomDelay()
	time.Sleep(delay)

	// Simulate random failures
	if m.shouldFail() {
		vm.ErrorMessage = "Failed to create VM: simulated error"
		return vm, errors.New(vm.ErrorMessage)
	}

	m.vms[id] = vm

	// Initialize VM properties
	// vm.Status = VMStatusRUNNING
	vm.CPU = cpu
	vm.Memory = memory
	vm.DiskUsage = 10.0 + rand.Float64()*30.0 // 10-40% initial disk usage

	return vm, nil
}

// StartVM starts a stopped VM
func (m *Manager) UpdateVM(id, name string, cpu int, memory int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	vm, exists := m.vms[id]
	if !exists {
		return fmt.Errorf("VM not found: %s", id)
	}

	if vm.Status != VMStatusSTOPPED && vm.Status != VMStatusSTARTED {
		return fmt.Errorf("VM cannot be updated from status %s", vm.Status)
	}

	delay := m.randomDelay()
	time.Sleep(delay)

	// Simulate random failures
	if m.shouldFail() {
		vm.ErrorMessage = "Failed to update VM: simulated error"
		return errors.New(vm.ErrorMessage)
	}

	// Initialize runtime properties
	vm.Name = name
	vm.CPU = cpu
	vm.Memory = memory
	vm.CPUUsage = 5.0 + rand.Float64()*45.0       // 5-50% initial CPU usage
	vm.MemoryUsage = 20.0 + rand.Float64()*40.0   // 20-60% initial memory usage
	vm.NetworkUsage = 50.0 + rand.Float64()*100.0 // 50-150 Mbps initial network throughput

	return nil
}

// StartVM starts a stopped VM
func (m *Manager) StartVM(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	vm, exists := m.vms[id]
	if !exists {
		return fmt.Errorf("VM not found: %s", id)
	}

	if vm.Status != VMStatusSTOPPED && vm.Status != VMStatusCREATED {
		return fmt.Errorf("VM cannot be started from status %s", vm.Status)
	}

	delay := m.randomDelay()
	time.Sleep(delay)

	// Simulate random failures
	if m.shouldFail() {
		vm.ErrorMessage = "Failed to start VM: simulated error"
		return errors.New(vm.ErrorMessage)
	}

	// Initialize runtime properties
	vm.Status = VMStatusSTARTED
	vm.CPUUsage = 5.0 + rand.Float64()*45.0       // 5-50% initial CPU usage
	vm.MemoryUsage = 20.0 + rand.Float64()*40.0   // 20-60% initial memory usage
	vm.NetworkUsage = 50.0 + rand.Float64()*100.0 // 50-150 Mbps initial network throughput

	return nil
}

// StopVM stops a running VM
func (m *Manager) StopVM(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	vm, exists := m.vms[id]
	if !exists {
		return fmt.Errorf("VM not found: %s", id)
	}

	if vm.Status != VMStatusSTARTED {
		return fmt.Errorf("VM cannot be stopped from status %s", vm.Status)
	}

	delay := m.randomDelay()
	time.Sleep(delay)

	// Simulate random failures
	if m.shouldFail() {
		vm.ErrorMessage = "Failed to stop VM: simulated error"
		return errors.New(vm.ErrorMessage)
	}

	// Update VM properties
	vm.Status = VMStatusSTOPPED
	vm.CPUUsage = 0
	vm.MemoryUsage = 0
	vm.NetworkUsage = 0

	return nil
}

// DeleteVM deletes a VM
func (m *Manager) DeleteVM(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	vm, exists := m.vms[id]
	if !exists {
		return fmt.Errorf("VM not found: %s", id)
	}

	if vm.Status != VMStatusSTOPPED {
		return fmt.Errorf("VM cannot be deleted from status %s", vm.Status)
	}

	delay := m.randomDelay()
	time.Sleep(delay)

	// Simulate random failures
	if m.shouldFail() {
		vm.ErrorMessage = "Failed to delete VM: simulated error"
		return errors.New(vm.ErrorMessage)
	}

	// Mark as deleted
	vm.Status = VMStatusDELETED

	return nil
}

// Retry
func (m *Manager) Retry(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	vm, exists := m.vms[id]
	if !exists {
		return fmt.Errorf("VM not found: %s", id)
	}

	// Not in error
	if vm.ErrorMessage == "" {
		return nil
	}

	delay := m.randomDelay()
	time.Sleep(delay)

	// Simulate random failures
	if m.shouldFail() {
		vm.ErrorMessage = "Failed to delete VM: simulated error"
		return errors.New(vm.ErrorMessage)
	}

	// Reset error
	vm.ErrorMessage = ""

	return nil
}

// UpdateVMResources periodically updates resource usage for running VMs
func (m *Manager) UpdateVMResources() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, vm := range m.vms {
		if vm.Status == VMStatusSTARTED {
			// Simulate changing resource utilization
			// CPU fluctuates more than memory
			vm.CPUUsage = clamp(vm.CPUUsage+(rand.Float64()*20.0-10.0), 1.0, 95.0)
			vm.MemoryUsage = clamp(vm.MemoryUsage+(rand.Float64()*10.0-3.0), 5.0, 90.0)
			vm.DiskUsage = clamp(vm.DiskUsage+(rand.Float64()*2.0-0.5), 10.0, 95.0)
			vm.NetworkUsage = clamp(vm.NetworkUsage+(rand.Float64()*30.0-15.0), 1.0, 500.0)
		}
	}
}

// GetStatusCounts returns the count of VMs in each status
func (m *Manager) GetStatusCounts() map[VMStatus]int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	counts := make(map[VMStatus]int)
	for _, vm := range m.vms {
		counts[vm.Status]++
	}

	return counts
}

// Helper methods
func (m *Manager) randomDelay() time.Duration {
	minDelay := m.delayRangeMin
	maxDelay := m.delayRangeMax

	// Calculate a random duration between min and max
	delta := maxDelay - minDelay
	if delta <= 0 {
		return minDelay
	}

	randomMs := rand.Int63n(int64(delta))
	return minDelay + time.Duration(randomMs)
}

func (m *Manager) shouldFail() bool {
	return rand.Float64() < m.errorRate
}

// Helper function to clamp a value between min and max
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
