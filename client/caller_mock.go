package client

import (
	"github.com/kuritsu/spyglass/api/types"
	"github.com/stretchr/testify/mock"
)

// CallerMock API client mock
type CallerMock struct {
	mock.Mock
}

// Init -ialize mock.
func (c *CallerMock) Init(url string) {
	c.Called(url)
}

// ListTargets operation
func (c *CallerMock) ListTargets(filter string, pageIndex int, pageSize int) ([]*types.Target, error) {
	args := c.Called(filter, pageIndex, pageSize)
	first := args.Get(0)
	if first != nil {
		return first.([]*types.Target), nil
	}
	return nil, args.Error(1)
}

// InsertOrUpdateMonitor operation.
func (c *CallerMock) InsertOrUpdateMonitor(monitor *types.Monitor) error {
	args := c.Called(monitor)
	return args.Error(0)
}

// InsertOrUpdateTarget operation.
func (c *CallerMock) InsertOrUpdateTarget(target *types.Target, forceStatusUpdate bool) error {
	args := c.Called(target, forceStatusUpdate)
	return args.Error(0)
}

// UpdateTargetStatus operation
func (c *CallerMock) UpdateTargetStatus(id string, status int, statusDescription string) error {
	args := c.Called(id, status, statusDescription)
	return args.Error(0)
}
