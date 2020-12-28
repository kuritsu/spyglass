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

func TestTargetGet(t *testing.T) {
	dbMock := testutil.StorageMock{
		GetTargetByIDResult: &types.Target{
			ID: "mytarget",
		},
	}
	r := Serve(&dbMock)
	w, jsonBytes := testutil.MakeRequest(http.MethodGet, "/targets/mytarget", nil, r)

	var target types.Target
	merr := json.Unmarshal(jsonBytes, &target)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, nil, merr)
	assert.Equal(t, "mytarget", target.ID)
}

func TestTargetGetNotFound(t *testing.T) {
	r := Serve(&testutil.StorageMock{})
	w, _ := testutil.MakeRequest(http.MethodGet, "/targets/mytarget", nil, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTargetGetDbError(t *testing.T) {
	dbMock := testutil.StorageMock{
		GetTargetByIDError: errors.New("Connection error"),
	}
	r := Serve(&dbMock)
	w, _ := testutil.MakeRequest(http.MethodGet, "/targets/mytarget", nil, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTargetPost(t *testing.T) {
	r := Serve(&testutil.StorageMock{})
	target := getValidTarget()
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/targets", target, r)

	assertValidTargetCreated(t, w, jsonBytes)
}

func TestTargetPostInvalidTarget(t *testing.T) {
	r := Serve(&testutil.StorageMock{})
	w, _ := testutil.MakeRequest(http.MethodPost, "/targets", "", r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTargetPostDbError(t *testing.T) {
	dbMock := testutil.StorageMock{
		InsertTargetError: errors.New("Connection error"),
	}
	r := Serve(&dbMock)
	target := getValidTarget()
	w, _ := testutil.MakeRequest(http.MethodPost, "/targets", target, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTargetPostErrorDuplicate(t *testing.T) {
	dbMock := testutil.StorageMock{
		InsertTargetError: errors.New("Duplicate"),
	}
	r := Serve(&dbMock)
	target := getValidTarget()
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/targets", target, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Duplicate target ID")
}

func TestTargetPostErrorInvalidID(t *testing.T) {
	r := Serve(&testutil.StorageMock{})
	target := getValidTarget()
	target.ID = ".target/"
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/targets", target, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Invalid target ID")
}

func TestTargetPostNoParentFoundError(t *testing.T) {
	r := Serve(&testutil.StorageMock{})
	target := getValidTarget()
	target.ID = "mytargets.target-1"
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/targets", target, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Target parent does not exist")
}

func TestTargetPostSearchingForParentError(t *testing.T) {
	r := Serve(&testutil.StorageMock{
		GetTargetByIDError: errors.New("Connection error"),
	})
	target := getValidTarget()
	target.ID = "mytargets.target-1"
	w, _ := testutil.MakeRequest(http.MethodPost, "/targets", target, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTargetPostParentExists(t *testing.T) {
	r := Serve(&testutil.StorageMock{
		GetTargetByIDResult: &types.Target{
			ID: "mytargets",
		},
	})
	target := getValidTarget()
	target.ID = "mytargets.target-1"
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/targets", target, r)

	assertValidTargetCreated(t, w, jsonBytes)
}

func TestTargetPostMonitorDoesntExist(t *testing.T) {
	r := Serve(&testutil.StorageMock{})
	target := getValidTarget()
	target.Monitor = &types.MonitorRef{MonitorID: "monitor1"}
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/targets", target, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Monitor does not exist")
}

func TestTargetPostGetMonitorHasDbError(t *testing.T) {
	r := Serve(&testutil.StorageMock{
		GetMonitorByIDError: errors.New("Connection error"),
	})
	target := getValidTarget()
	target.Monitor = &types.MonitorRef{MonitorID: "monitor1"}
	w, _ := testutil.MakeRequest(http.MethodPost, "/targets", target, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTargetPostMonitorExists(t *testing.T) {
	r := Serve(&testutil.StorageMock{
		GetMonitorByIDResult: &types.Monitor{ID: "monitor1"},
	})
	target := getValidTarget()
	target.Monitor = &types.MonitorRef{MonitorID: "monitor1"}
	w, jsonBytes := testutil.MakeRequest(http.MethodPost, "/targets", target, r)

	assertValidTargetCreated(t, w, jsonBytes)
}

func assertValidTargetCreated(t *testing.T, w *httptest.ResponseRecorder, jsonBytes []byte) {
	var newTarget types.Target
	merr := json.Unmarshal(jsonBytes, &newTarget)

	assert.Equal(t, nil, merr)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEqual(t, time.Time{}, newTarget.CreatedAt)
}

func getValidTarget() *types.Target {
	return &types.Target{
		ID:          "mytarget-1",
		Description: "my target",
	}
}
