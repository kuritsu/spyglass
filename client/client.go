package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kuritsu/spyglass/api/types"
	logr "github.com/sirupsen/logrus"
)

// APIClient for calling Spyglass API.
type APIClient struct {
	client *http.Client
	log    *logr.Logger
	url    string
}

// Create API client.
func Create(log *logr.Logger) APICaller {
	return &APIClient{
		client: &http.Client{},
		log:    log}
}

// Init -ialize the client with the url.
func (c *APIClient) Init(url string) {
	c.url = url
}

// ListTargets operation
func (c *APIClient) ListTargets(filter string, pageIndex int, pageSize int) ([]*types.Target, error) {
	c.log.Debugf("Getting targets %v, pageIndex %v, pageSize %v", filter, pageIndex, pageSize)
	response, err := c.client.Get(
		fmt.Sprintf("%s/targets?filter=%s&pageIndex=%v&pageSize=%v",
			c.url, filter, pageIndex, pageSize))
	if err != nil {
		return nil, err
	}
	result := []*types.Target{}
	bodyBytes, rerr := ioutil.ReadAll(response.Body)
	if rerr != nil {
		return nil, rerr
	}
	if err = json.Unmarshal(bodyBytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

// InsertOrUpdateMonitor operation.
func (c *APIClient) InsertOrUpdateMonitor(monitor *types.Monitor) error {
	c.log.Debug("Getting monitor ", monitor.ID)
	response, err := c.client.Get(fmt.Sprintf("%s/monitors/%s", c.url, monitor.ID))
	if err != nil {
		return err
	}
	bodyBytes, _ := json.Marshal(monitor)
	reader := strings.NewReader(string(bodyBytes))
	var request *http.Request
	switch response.StatusCode {
	case http.StatusNotFound:
		c.log.Debug("Posting monitor ", monitor.ID)
		response, err = c.client.Post(fmt.Sprintf("%s/monitors", c.url), "application/json", reader)
	case http.StatusOK:
		c.log.Debug("Putting monitor ", monitor.ID)
		request, _ = http.NewRequest("PUT", fmt.Sprintf("%s/monitors/%s", c.url, monitor.ID), reader)
		request.Header["Content-Type"] = []string{"application/json"}
		response, err = c.client.Do(request)
	}
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated &&
		response.StatusCode != http.StatusOK {
		return c.errorWithMessage(response)
	}
	return nil
}

// InsertOrUpdateTarget operation.
func (c *APIClient) InsertOrUpdateTarget(target *types.Target, forceStatusUpdate bool) error {
	c.log.Debug("Getting target ", target.ID)
	response, err := c.client.Get(fmt.Sprintf("%s/targets/%s", c.url, target.ID))
	if err != nil {
		return err
	}
	bodyBytes, _ := json.Marshal(target)
	reader := strings.NewReader(string(bodyBytes))
	var request *http.Request
	switch response.StatusCode {
	case http.StatusNotFound:
		c.log.Debug("Posting target ", target.ID)
		response, err = c.client.Post(fmt.Sprintf("%s/targets", c.url), "application/json", reader)
	case http.StatusOK:
		c.log.Debug("Putting target ", target.ID)
		request, _ = http.NewRequest(http.MethodPut,
			fmt.Sprintf("%s/targets/%s?forceStatusUpdate=%v", c.url, target.ID, forceStatusUpdate), reader)
		request.Header["Content-Type"] = []string{"application/json"}
		response, err = c.client.Do(request)
	}
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
		fmt.Sprintf("%s/targets/%s", c.url, id), reader)
	request.Header["Content-Type"] = []string{"application/json"}
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

func (c *APIClient) errorWithMessage(response *http.Response) error {
	responseBytes, _ := ioutil.ReadAll(response.Body)
	errorMsg := make(map[string]string)
	json.Unmarshal(responseBytes, &errorMsg)
	return errors.New(errorMsg["message"])
}
