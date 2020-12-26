package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kuritsu/spyglass/api/storage/testutil"
	"github.com/kuritsu/spyglass/api/types"
	"github.com/stretchr/testify/assert"
)

func TestMonitorPost(t *testing.T) {
	r := Serve(&testutil.Mock{})
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

	jsonBytes, _ := ioutil.ReadAll(w.Result().Body)
	var newMonitor types.Monitor
	fmt.Println(string(jsonBytes))
	json.Unmarshal(jsonBytes, &newMonitor)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEqual(t, time.Time{}, newMonitor.CreatedAt)
}

func TestMonitorPostInvalidMonitor(t *testing.T) {
	r := Serve(&testutil.Mock{})
	w := httptest.NewRecorder()
	jsonBody := "{}"
	req, _ := http.NewRequest(http.MethodPost, "/monitors",
		strings.NewReader(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMonitorPostDbError(t *testing.T) {
	dbMock := testutil.Mock{
		InsertMonitorError: errors.New("Connection error"),
	}
	r := Serve(&dbMock)
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
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMonitorGet(t *testing.T) {
	dbMock := testutil.Mock{
		GetMonitorByIDResult: &types.Monitor{
			ID: "mymonitor",
		},
	}
	r := Serve(&dbMock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/monitors/mymonitor", nil)
	r.ServeHTTP(w, req)

	jsonBytes, _ := ioutil.ReadAll(w.Result().Body)
	var monitor types.Monitor
	fmt.Println(string(jsonBytes))
	merr := json.Unmarshal(jsonBytes, &monitor)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, nil, merr)
	assert.Equal(t, "mymonitor", monitor.ID)
}

func TestMonitorGetNotFound(t *testing.T) {
	dbMock := testutil.Mock{}
	r := Serve(&dbMock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/monitors/mymonitor", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMonitorGetDbError(t *testing.T) {
	dbMock := testutil.Mock{
		GetMonitorByIDError: errors.New("Connection error"),
	}
	r := Serve(&dbMock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/monitors/mymonitor", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
