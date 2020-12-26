package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/api/types"
	"github.com/stretchr/testify/assert"
)

func TestMonitorPost(t *testing.T) {
	r := Serve(&storage.Mock{})
	w := httptest.NewRecorder()
	monitor := types.Monitor{
		ID:       "1",
		Type:     "docker",
		Schedule: "* * * * *",
		Definition: types.MonitorDefinition{
			DockerDefinition: types.DockerDefinition{
				Image: "nginx:latest",
				DockerEnv: map[string]string{
					"val1": "val2",
				},
			},
		},
	}
	jsonBody, _ := json.Marshal(monitor)
	req, _ := http.NewRequest(http.MethodPost, "/monitors",
		strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
