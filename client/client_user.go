package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kuritsu/spyglass/api/types"
)

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

func (c *APIClient) ListUsers(pageIndex, pageSize int) ([]*types.User, error) {
	c.log.Debug("Getting users...")
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users?pageIndex=%v&pageSize=%v", c.url, pageIndex, pageSize), http.NoBody)
	c.addAuthHeader(req)
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	result := []*types.User{}
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

func (c *APIClient) CreateUserToken(email string, expiration time.Time) (string, error) {
	c.log.Debug("CreateUserToken ", email)
	userTokenReq := types.UserTokenRequest{
		Expiration: expiration,
	}
	bodyBytes, _ := json.Marshal(userTokenReq)
	reader := strings.NewReader(string(bodyBytes))
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/user/token/%s", c.url, email), reader)
	c.addAuthHeader(request)
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
	c.log.Debug("Create user token success.")
	return result, nil
}
