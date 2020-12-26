package storage

import (
	"github.com/kuritsu/spyglass/api/types"
	"github.com/stretchr/testify/mock"
)

// Mock is a mock for storage
type Mock struct {
	mock.Mock
}

// Initialize with nothing
func (m *Mock) Initialize() {}

// GetMonitorByID returns nothing
func (m *Mock) GetMonitorByID(id string) *types.Monitor {
	return nil
}

// GetTargetByID returns nothing
func (m *Mock) GetTargetByID(id string) *types.Target {
	return nil
}
