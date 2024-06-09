package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/kuritsu/spyglass/api/types"
)

// ListTargets operation
func (c *APIClient) ListTargets(filter string, pageIndex int, pageSize int) ([]*types.Target, error) {
	c.log.Debugf("Getting targets %v, pageIndex %v, pageSize %v", filter, pageIndex, pageSize)
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/targets?contains=%s&pageIndex=%v&pageSize=%v",
		c.url, filter, pageIndex, pageSize), http.NoBody)
	c.addAuthHeader(req)
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	result := []*types.Target{}
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

// InsertOrUpdateTarget operation.
func (c *APIClient) InsertOrUpdateTarget(target *types.Target, forceStatusUpdate bool) error {
	c.log.Debug("Getting target ", target.ID)
	var request *http.Request
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/targets/%s", c.url, target.ID), http.NoBody)
	c.addAuthHeader(request)
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	bodyBytes, _ := json.Marshal(target)
	reader := strings.NewReader(string(bodyBytes))
	switch response.StatusCode {
	case http.StatusNotFound:
		c.log.Debug("Posting target ", target.ID)
		request, _ = http.NewRequest(http.MethodPost,
			fmt.Sprintf("%s/targets", c.url), reader)
	case http.StatusOK:
		c.log.Debug("Putting target ", target.ID)
		request, _ = http.NewRequest(http.MethodPut,
			fmt.Sprintf("%s/targets/%s?forceStatusUpdate=%v", c.url, target.ID, forceStatusUpdate), reader)
	}
	c.addAuthHeader(request)
	response, err = c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated &&
		response.StatusCode != http.StatusOK {
		return c.errorWithMessage(response)
	}
	return nil
}

// UpdateTargetStatus operation
func (c *APIClient) UpdateTargetStatus(id string, status int, statusDescription string) error {
	targetPatch := make(map[string]interface{})
	targetPatch["status"] = status
	if statusDescription != "" {
		targetPatch["statusDescription"] = statusDescription
	}
	bodyBytes, _ := json.Marshal(targetPatch)
	reader := strings.NewReader(string(bodyBytes))
	c.log.Debug("Patching target ", id)
	request, _ := http.NewRequest(http.MethodPatch,
		fmt.Sprintf("%s/target?id=%s", c.url, id), reader)
	c.addAuthHeader(request)
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return c.errorWithMessage(response)
	}
	c.log.Debug("Status patched successfully.")
	return nil
}

// UpdateTargetStatus operation
func (c *APIClient) GetTargetByID(id string, includeChildren bool) (types.TargetRef, error) {
	c.log.Debug("Getting target ", id)
	request, _ := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/target?id=%s&includeChildren=%v", c.url, url.QueryEscape(id), includeChildren), http.NoBody)
	c.addAuthHeader(request)
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, c.errorWithMessage(response)
	}
	result := &types.Target{}
	bodyBytes, rerr := io.ReadAll(response.Body)
	fmt.Printf("%v\n", string(bodyBytes))
	if rerr != nil {
		return nil, rerr
	}
	if err = json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}
	c.log.Debug("Get target successfully.")
	return result, nil
}

func (c *APIClient) DeleteTarget(id string) (int, error) {
	c.log.Debug("Delete target ", id)
	request, _ := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/target?id=%s", c.url, url.QueryEscape(id)), http.NoBody)
	c.addAuthHeader(request)
	response, err := c.client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, c.errorWithMessage(response)
	}
	responseBytes, _ := io.ReadAll(response.Body)
	responseMsg := make(map[string]int)
	json.Unmarshal(responseBytes, &responseMsg)
	deleteCount := responseMsg["deletedCount"]
	return deleteCount, nil
}
