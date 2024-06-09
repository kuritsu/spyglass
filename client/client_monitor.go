package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kuritsu/spyglass/api/types"
)

// InsertOrUpdateMonitor operation.
func (c *APIClient) InsertOrUpdateMonitor(monitor *types.Monitor) error {
	c.log.Debug("Getting monitor ", monitor.ID)
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/monitors/%s", c.url, monitor.ID), http.NoBody)
	c.addAuthHeader(request)
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	bodyBytes, _ := json.Marshal(monitor)
	reader := strings.NewReader(string(bodyBytes))
	switch response.StatusCode {
	case http.StatusNotFound:
		c.log.Debug("Posting monitor ", monitor.ID)
		request, _ = http.NewRequest("POST", fmt.Sprintf("%s/monitors", c.url), reader)
	case http.StatusOK:
		c.log.Debug("Putting monitor ", monitor.ID)
		request, _ = http.NewRequest("PUT", fmt.Sprintf("%s/monitors/%s", c.url, monitor.ID), reader)
	}
	c.addAuthHeader(request)
	response, err = c.client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated &&
		response.StatusCode != http.StatusOK {
		return c.errorWithMessage(response)
	}
	return nil
}

// ListMonitors operation
func (c *APIClient) ListMonitors(pageIndex int, pageSize int) ([]*types.Monitor, error) {
	c.log.Debugf("Getting monitors, pageIndex %v, pageSize %v", pageIndex, pageSize)
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/monitors?pageIndex=%v&pageSize=%v",
		c.url, pageIndex, pageSize), http.NoBody)
	c.addAuthHeader(req)
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	result := []*types.Monitor{}
	bodyBytes, rerr := io.ReadAll(response.Body)
	if rerr != nil {
		return nil, rerr
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(bodyBytes))
	}
	if err = json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}
