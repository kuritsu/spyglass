package testutil

import (
	"time"

	"github.com/kuritsu/spyglass/api/types"
	"github.com/stretchr/testify/mock"
)

// Mock is a mock for storage
type Mock struct {
	mock.Mock
	GetMonitorByIDResult *types.Monitor
	GetMonitorByIDError  error
	GetTargetByIDResult  *types.Target
	GetTargetByIDError   error
	InsertMonitorError   error
	InsertTargetError    error
}

// Init with nothing
func (m *Mock) Init() {}

// Free resources
func (m *Mock) Free() {}

// GetMonitorByID returns nothing
func (m *Mock) GetMonitorByID(id string) (*types.Monitor, error) {
	return m.GetMonitorByIDResult, m.GetMonitorByIDError
}

// InsertMonitor in the db
func (m *Mock) InsertMonitor(monitor *types.Monitor) (*types.Monitor, error) {
	monitor.CreatedAt = time.Now()
	return monitor, m.InsertMonitorError
}

// GetTargetByID returns nothing
func (m *Mock) GetTargetByID(id string) (*types.Target, error) {
	return m.GetTargetByIDResult, m.GetTargetByIDError
}

// InsertTarget into the db.
func (m *Mock) InsertTarget(target *types.Target) (*types.Target, error) {
	target.CreatedAt = time.Now()
	return target, m.InsertTargetError
}
