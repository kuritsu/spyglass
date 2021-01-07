package client

import (
	"net/http"

	"github.com/kuritsu/spyglass/api/types"
	logr "github.com/sirupsen/logrus"
)

// APIClient for calling Spyglass API.
type APIClient struct {
	client *http.Client
	log    *logr.Logger
}

// InsertOrUpdateMonitor operation.
func (c *APIClient) InsertOrUpdateMonitor(monitor *types.Monitor) error {
	c.log.Debug("Getting monitor", monitor.ID)
	return nil
}

// InsertOrUpdateTarget operation.
func (c *APIClient) InsertOrUpdateTarget(target *types.Target) error {
	c.log.Debug("Getting target", target.ID)
	return nil
}

// Create API client.
func Create(log *logr.Logger) APICaller {
	return &APIClient{&http.Client{}, log}
}
