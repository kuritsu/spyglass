package client

import (
	"net/http"

	"github.com/kuritsu/spyglass/api/types"
)

// APIClient for calling Spyglass API.
type APIClient struct {
	client *http.Client
}

// InsertOrUpdateMonitor operation.
func (c *APIClient) InsertOrUpdateMonitor(monitor *types.Monitor) error {
	return nil
}

// InsertOrUpdateTarget operation.
func (c *APIClient) InsertOrUpdateTarget(target *types.Target) error {
	return nil
}
