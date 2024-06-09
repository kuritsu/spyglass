package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

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
	if token, ok := os.LookupEnv("SPYGLASS_TOKEN"); ok {
		c.token = token
		return
	}
	homedir, _ := os.UserHomeDir()
	fname := filepath.Join(homedir, ".spyglass.token")
	tokenBytes, err := os.ReadFile(fname)
	if err == nil {
		c.token = string(tokenBytes)
	}
}

func (c *APIClient) errorWithMessage(response *http.Response) error {
	responseBytes, _ := io.ReadAll(response.Body)
	errorMsg := make(map[string]string)
	json.Unmarshal(responseBytes, &errorMsg)
	return errors.New(errorMsg["message"])
}

func (c *APIClient) addAuthHeader(request *http.Request) {
	request.Header["Content-Type"] = []string{"application/json"}
	request.Header.Add("Authorization", c.token)
}
