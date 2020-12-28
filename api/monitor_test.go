package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kuritsu/spyglass/api/testutil"
	"github.com/kuritsu/spyglass/api/types"
	"github.com/stretchr/testify/assert"
)

func TestMonitorGet(t *testing.T) {
	dbMock := testutil.StorageMock{
		GetMonitorByIDResult: &types.Monitor{
			ID: "mymonitor",
		},
	}
	r := Serve(&dbMock)
	w, jsonBytes := testutil.MakeRequest(http.MethodGet, "/monitors/mymonitor", nil, r)

	var monitor types.Monitor
	merr := json.Unmarshal(jsonBytes, &monitor)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, nil, merr)
	assert.Equal(t, "mymonitor", monitor.ID)
}

func TestMonitorGetNotFound(t *testing.T) {
	dbMock := testutil.StorageMock{}
	r := Serve(&dbMock)
	w, _ := testutil.MakeRequest(http.MethodGet, "/monitors/mymonitor", nil, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMonitorGetDbError(t *testing.T) {
	dbMock := testutil.StorageMock{
		GetMonitorByIDError: errors.New("Connection error"),
	}
	r := Serve(&dbMock)
	w, _ := testutil.MakeRequest(http.MethodGet, "/monitors/mymonitor", nil, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMonitorPost(t *testing.T) {
	r := Serve(&testutil.StorageMock{})
	monitor := getValidMonitor()
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/monitors", monitor, r)

	assertValidMonitorCreated(t, w, jsonBytes)
}

func TestMonitorPostInvalidMonitor(t *testing.T) {
	r := Serve(&testutil.StorageMock{})
	w, _ := testutil.MakeRequest(http.MethodPost, "/monitors", "", r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMonitorPostDbError(t *testing.T) {
	dbMock := testutil.StorageMock{
		InsertMonitorError: errors.New("Connection error"),
	}
	r := Serve(&dbMock)
	monitor := getValidMonitor()
	w, _ := testutil.MakeRequest(http.MethodPost, "/monitors", monitor, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMonitorPostErrorDuplicate(t *testing.T) {
	dbMock := testutil.StorageMock{
		InsertMonitorError: errors.New("Duplicate"),
	}
	r := Serve(&dbMock)
	monitor := getValidMonitor()
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/monitors", monitor, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Duplicate monitor ID")
}

func TestMonitorPostErrorInvalidID(t *testing.T) {
	r := Serve(&testutil.StorageMock{})
	monitor := getValidMonitor()
	monitor.ID = "/monitor"
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/monitors", monitor, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Invalid monitor ID")
}

func assertValidMonitorCreated(t *testing.T, w *httptest.ResponseRecorder, jsonBytes []byte) {
	var newMonitor types.Monitor
	merr := json.Unmarshal(jsonBytes, &newMonitor)

	assert.Equal(t, nil, merr)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEqual(t, time.Time{}, newMonitor.CreatedAt)
}

func getValidMonitor() *types.Monitor {
	return &types.Monitor{
		ID:       "monitors.mymonitor-1",
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
}
