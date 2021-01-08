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

// Init -ialize the client with the url.
func (c *APIClient) Init(url string) {
	c.url = url
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
		request, err = http.NewRequest("PUT", fmt.Sprintf("%s/monitors/%s", c.url, monitor.ID), reader)
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

// Create API client.
func Create(log *logr.Logger) APICaller {
	return &APIClient{
		client: &http.Client{},
		log:    log}
}

func (c *APIClient) errorWithMessage(response *http.Response) error {
	responseBytes, _ := ioutil.ReadAll(response.Body)
	errorMsg := make(map[string]string)
	json.Unmarshal(responseBytes, &errorMsg)
	return errors.New(errorMsg["message"])
}
