package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kuritsu/spyglass/api/types"
	logr "github.com/sirupsen/logrus"
)

// APIClient for calling Spyglass API.
type APIClient struct {
	client *http.Client
	log    *logr.Logger
	url    string
	token  string
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
	homedir, _ := os.UserHomeDir()
	fname := filepath.Join(homedir, ".spyglass.token")
	tokenBytes, err := os.ReadFile(fname)
	if err == nil {
		c.token = string(tokenBytes)
	}
}

// ListTargets operation
func (c *APIClient) ListTargets(filter string, pageIndex int, pageSize int) ([]*types.Target, error) {
	c.log.Debugf("Getting targets %v, pageIndex %v, pageSize %v", filter, pageIndex, pageSize)
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/targets?contains=%s&pageIndex=%v&pageSize=%v",
		c.url, filter, pageIndex, pageSize), http.NoBody)
	req.Header.Add("Authorization", c.token)
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

// InsertOrUpdateMonitor operation.
func (c *APIClient) InsertOrUpdateMonitor(monitor *types.Monitor) error {
	var request *http.Request
	c.log.Debug("Getting monitor ", monitor.ID)
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/monitors/%s", c.url, monitor.ID), http.NoBody)
	request.Header.Add("Authorization", c.token)
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
	request.Header["Content-Type"] = []string{"application/json"}
	request.Header.Add("Authorization", c.token)
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

// InsertOrUpdateTarget operation.
func (c *APIClient) InsertOrUpdateTarget(target *types.Target, forceStatusUpdate bool) error {
	c.log.Debug("Getting target ", target.ID)
	var request *http.Request
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/targets/%s", c.url, target.ID), http.NoBody)
	request.Header.Add("Authorization", c.token)
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
	request.Header["Content-Type"] = []string{"application/json"}
	request.Header.Add("Authorization", c.token)
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
	request.Header["Content-Type"] = []string{"application/json"}
	request.Header.Add("Authorization", c.token)
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
	request.Header["Content-Type"] = []string{"application/json"}
	request.Header.Add("Authorization", c.token)
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

func (c *APIClient) errorWithMessage(response *http.Response) error {
	responseBytes, _ := io.ReadAll(response.Body)
	errorMsg := make(map[string]string)
	json.Unmarshal(responseBytes, &errorMsg)
	return errors.New(errorMsg["message"])
}

func (c *APIClient) Login(email string, password string) (string, error) {
	c.log.Debug("Login ", email)
	userPwdDict := types.AuthRequest{
		Email:    email,
		Password: password,
	}
	bodyBytes, _ := json.Marshal(userPwdDict)
	reader := strings.NewReader(string(bodyBytes))
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/user/login", c.url), reader)
	request.Header["Content-Type"] = []string{"application/json"}
	response, err := c.client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", c.errorWithMessage(response)
	}
	var result string
	bodyBytes, rerr := io.ReadAll(response.Body)
	if rerr != nil {
		return "", rerr
	}
	if err = json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}
	c.log.Debug("Auth success.")
	return result, nil
}
