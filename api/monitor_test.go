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
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMonitorGet(t *testing.T) {
	dbMock := testutil.StorageMock{
		GetMonitorByIDResult: &types.Monitor{
			ID: "mymonitor",
		},
	}
	r := Create(&dbMock, logrus.New()).Serve()
	w, jsonBytes := testutil.MakeRequest(http.MethodGet, "/monitors/mymonitor", nil, r)

	var monitor types.Monitor
	merr := json.Unmarshal(jsonBytes, &monitor)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, nil, merr)
	assert.Equal(t, "mymonitor", monitor.ID)
}

func TestMonitorGetNotFound(t *testing.T) {
	dbMock := testutil.StorageMock{}
	r := Create(&dbMock, logrus.New()).Serve()
	w, _ := testutil.MakeRequest(http.MethodGet, "/monitors/mymonitor", nil, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMonitorGetDbError(t *testing.T) {
	dbMock := testutil.StorageMock{
		GetMonitorByIDError: errors.New("Connection error"),
	}
	r := Create(&dbMock, logrus.New()).Serve()
	w, _ := testutil.MakeRequest(http.MethodGet, "/monitors/mymonitor", nil, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMonitorGetAll(t *testing.T) {
	dbMock := testutil.StorageMock{
		GetAllMonitorsResult: []types.Monitor{
			{ID: "1"},
			{ID: "2"},
			{ID: "3"},
		},
	}
	r := Create(&dbMock, logrus.New()).Serve()
	w, jsonBytes := testutil.MakeRequest(http.MethodGet, "/monitors", nil, r)

	var monitors []types.Monitor
	merr := json.Unmarshal(jsonBytes, &monitors)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, nil, merr)
	assert.Equal(t, 3, len(monitors))
}

func TestMonitorGetAllWithPageSize(t *testing.T) {
	dbMock := testutil.StorageMock{
		GetAllMonitorsResult: []types.Monitor{
			{ID: "2"},
		},
	}
	r := Create(&dbMock, logrus.New()).Serve()
	w, jsonBytes := testutil.MakeRequest(http.MethodGet, "/monitors?pageSize=1&pageIndex=0", nil, r)

	var monitors []types.Monitor
	json.Unmarshal(jsonBytes, &monitors)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, len(monitors))
}

func TestMonitorGetAllEmptyList(t *testing.T) {
	dbMock := testutil.StorageMock{
		GetAllMonitorsResult: []types.Monitor{},
	}
	r := Create(&dbMock, logrus.New()).Serve()
	w, jsonBytes := testutil.MakeRequest(http.MethodGet, "/monitors", nil, r)

	var monitors []types.Monitor
	merr := json.Unmarshal(jsonBytes, &monitors)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, merr, nil)
	assert.NotNil(t, monitors)
	assert.Equal(t, len(monitors), 0)
}

func TestMonitorGetAllWithDbError(t *testing.T) {
	dbMock := testutil.StorageMock{
		GetAllMonitorsError: errors.New("Connection error"),
	}
	r := Create(&dbMock, logrus.New()).Serve()
	w, _ := testutil.MakeRequest(http.MethodGet, "/monitors", nil, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMonitorGetAllWithInvalidPageSize(t *testing.T) {
	r := Create(&testutil.StorageMock{}, logrus.New()).Serve()
	w, jsonBytes := testutil.MakeRequest(http.MethodGet, "/monitors?pageSize=asd", nil, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Invalid page size")
}

func TestMonitorGetAllWithInvalidPageIndex(t *testing.T) {
	r := Create(&testutil.StorageMock{}, logrus.New()).Serve()
	w, jsonBytes := testutil.MakeRequest(http.MethodGet, "/monitors?pageIndex=asd", nil, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Invalid page index")
}

func TestMonitorGetAllWithInvalidContains(t *testing.T) {
	r := Create(&testutil.StorageMock{}, logrus.New()).Serve()
	w, jsonBytes := testutil.MakeRequest(http.MethodGet, "/monitors?contains=a,sd", nil, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Invalid contains expression")
}

func TestMonitorPost(t *testing.T) {
	r := Create(&testutil.StorageMock{}, logrus.New()).Serve()
	monitor := getValidMonitor()
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/monitors", monitor, r)

	assertValidMonitorCreated(t, w, jsonBytes)
}

func TestMonitorPostInvalidMonitor(t *testing.T) {
	r := Create(&testutil.StorageMock{}, logrus.New()).Serve()
	w, _ := testutil.MakeRequest(http.MethodPost, "/monitors", "", r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMonitorPostDbError(t *testing.T) {
	dbMock := testutil.StorageMock{
		InsertMonitorError: errors.New("Connection error"),
	}
	r := Create(&dbMock, logrus.New()).Serve()
	monitor := getValidMonitor()
	w, _ := testutil.MakeRequest(http.MethodPost, "/monitors", monitor, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMonitorPostErrorDuplicate(t *testing.T) {
	dbMock := testutil.StorageMock{
		InsertMonitorError: errors.New("Duplicate"),
	}
	r := Create(&dbMock, logrus.New()).Serve()
	monitor := getValidMonitor()
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/monitors", monitor, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Duplicate monitor ID")
}

func TestMonitorPostErrorInvalidID(t *testing.T) {
	r := Create(&testutil.StorageMock{}, logrus.New()).Serve()
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
	assert.Equal(t, w.Code, http.StatusCreated)
	assert.NotEqual(t, newMonitor.Permissions.CreatedAt, time.Time{})
}

func getValidMonitor() *types.Monitor {
	return &types.Monitor{
		ID:       "monitors.mymonitor-1",
		Type:     "docker",
		Schedule: "* * * * *",
		Definition: &types.MonitorDefinition{
			Docker: &types.DockerDefinition{
				Image: "nginx:latest",
				DockerEnv: map[string]string{
					"val1": "val2",
				},
			},
		},
	}
}
