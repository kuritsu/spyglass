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
func (c *CallerMock) Init(url string) {}

// InsertOrUpdateMonitor operation.
func (c *CallerMock) InsertOrUpdateMonitor(monitor *types.Monitor) error {
	args := c.Called(monitor)
	return args.Error(0)
}

// InsertOrUpdateTarget operation.
func (c *CallerMock) InsertOrUpdateTarget(target *types.Target) error {
	args := c.Called(target)
	return args.Error(0)
}
