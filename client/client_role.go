package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kuritsu/spyglass/api/types"
)

func (c *APIClient) InsertRole(role *types.Role) error {
	c.log.Debug("Add role ", role.Name)
	bodyBytes, _ := json.Marshal(role)
	reader := strings.NewReader(string(bodyBytes))
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/roles", c.url), reader)
	c.addAuthHeader(request)
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return c.errorWithMessage(response)
	}
	c.log.Debug("Add role success.")
	return nil
}

func (c *APIClient) UpdateRole(role string, usersAdd, usersRemove []string) error {
	c.log.Debug("Update role ", role)
	roleUpdateReq := types.RoleUpdateRequest{
		UsersAdd:    usersAdd,
		UsersRemove: usersRemove,
	}
	bodyBytes, _ := json.Marshal(roleUpdateReq)
	reader := strings.NewReader(string(bodyBytes))
	request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/role/%s", c.url, role), reader)
	c.addAuthHeader(request)
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return c.errorWithMessage(response)
	}
	c.log.Debug("Update role success.")
	return nil
}

func (c *APIClient) ListRoles(pageIndex, pageSize int) ([]*types.Role, error) {
	c.log.Debug("Getting roles...")
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/roles?pageIndex=%v&pageSize=%v", c.url, pageIndex, pageSize), http.NoBody)
	c.addAuthHeader(req)
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	result := []*types.Role{}
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
