package testutil

import (
	"time"

	"github.com/kuritsu/spyglass/api/types"
	"github.com/stretchr/testify/mock"
)

// Mock is a mock for storage
type Mock struct {
	mock.Mock
	InsertMonitorError error
}

// Init with nothing
func (m *Mock) Init() {}

// Free resources
func (m *Mock) Free() {}

// GetMonitorByID returns nothing
func (m *Mock) GetMonitorByID(id string) (*types.Monitor, error) {
	return nil, nil
}

// InsertMonitor in the db
func (m *Mock) InsertMonitor(monitor *types.Monitor) (*types.Monitor, error) {
	monitor.CreatedAt = time.Now()
	return monitor, m.InsertMonitorError
}

// GetTargetByID returns nothing
func (m *Mock) GetTargetByID(id string) (*types.Target, error) {
	return nil, nil
}
