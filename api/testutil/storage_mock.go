package testutil

import (
	"time"

	"github.com/kuritsu/spyglass/api/types"
	"github.com/stretchr/testify/mock"
)

// StorageMock is a mock for storage
type StorageMock struct {
	mock.Mock
	// TODO: This needs to be migrated to testify/mock way
	GetMonitorByIDResult    *types.Monitor
	GetMonitorByIDError     error
	GetAllMonitorsResult    []types.Monitor
	GetAllMonitorsError     error
	GetTargetByIDResult     *types.Target
	GetTargetByIDError      error
	GetAllTargetsResult     []*types.Target
	GetAllTargetsError      error
	InsertMonitorError      error
	InsertTargetError       error
	UpdateTargetStatusError error
}

// Init with nothing
func (m *StorageMock) Init() {}

// Free resources
func (m *StorageMock) Free() {}

// GetMonitorByID returns mocked stuff
func (m *StorageMock) GetMonitorByID(id string) (*types.Monitor, error) {
	return m.GetMonitorByIDResult, m.GetMonitorByIDError
}

// GetAllMonitors returns mocked stuff
func (m *StorageMock) GetAllMonitors(pageSize int64, pageIndex int64, contains string) ([]types.Monitor, error) {
	return m.GetAllMonitorsResult, m.GetAllMonitorsError
}

// GetAllTargets returns mocked stuff
func (m *StorageMock) GetAllTargets(pageSize int64, pageIndex int64, contains string) ([]*types.Target, error) {
	return m.GetAllTargetsResult, m.GetAllTargetsError
}

// InsertMonitor in the db
func (m *StorageMock) InsertMonitor(monitor *types.Monitor) (*types.Monitor, error) {
	monitor.CreatedAt = time.Now()
	return monitor, m.InsertMonitorError
}

// GetTargetByID returns nothing
func (m *StorageMock) GetTargetByID(id string, includeChildren bool) (*types.Target, error) {
	return m.GetTargetByIDResult, m.GetTargetByIDError
}

// InsertTarget into the db.
func (m *StorageMock) InsertTarget(target *types.Target) (*types.Target, error) {
	target.CreatedAt = time.Now()
	return target, m.InsertTargetError
}

// UpdateTargetStatus with all modified fields
func (m *StorageMock) UpdateTargetStatus(target *types.Target, targetPatch *types.TargetPatch) (*types.Target, error) {
	target.Status = targetPatch.Status
	if targetPatch.StatusDescription != "" {
		target.StatusDescription = targetPatch.StatusDescription
	}
	return target, m.UpdateTargetStatusError
}

// UpdateTarget with also a status update force flag.
func (m *StorageMock) UpdateTarget(oldTarget *types.Target, newTarget *types.Target,
	forceStatusUpdate bool) (*types.Target, error) {
	args := m.Called(oldTarget, newTarget, forceStatusUpdate)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*types.Target), nil
}
