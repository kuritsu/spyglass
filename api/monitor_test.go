package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitorPost(t *testing.T) {
	r := Serve()
	w := httptest.NewRecorder()
	monitor := Monitor{
		ID:       "1",
		Type:     "docker",
		Schedule: "* * * * *",
		Definition: MonitorDefinition{
			DockerDefinition: DockerDefinition{
				Image: "nginx:latest",
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
