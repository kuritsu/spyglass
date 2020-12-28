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

func TestTargetGet(t *testing.T) {
	dbMock := testutil.Mock{
		GetTargetByIDResult: &types.Target{
			ID: "mytarget",
		},
	}
	r := Serve(&dbMock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/targets/mytarget", nil)
	r.ServeHTTP(w, req)

	jsonBytes, _ := ioutil.ReadAll(w.Result().Body)
	var monitor types.Monitor
	fmt.Println(string(jsonBytes))
	merr := json.Unmarshal(jsonBytes, &monitor)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, nil, merr)
	assert.Equal(t, "mytarget", monitor.ID)
}

func TestTargetGetNotFound(t *testing.T) {
	dbMock := testutil.Mock{}
	r := Serve(&dbMock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/targets/mytarget", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTargetGetDbError(t *testing.T) {
	dbMock := testutil.Mock{
		GetTargetByIDError: errors.New("Connection error"),
	}
	r := Serve(&dbMock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/targets/mytarget", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTargetPost(t *testing.T) {
	r := Serve(&testutil.Mock{})
	w := httptest.NewRecorder()
	target := getValidTarget()
	jsonBody, _ := json.Marshal(target)
	req, _ := http.NewRequest(http.MethodPost, "/targets",
		strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	jsonBytes, _ := ioutil.ReadAll(w.Result().Body)
	var newTarget types.Target
	json.Unmarshal(jsonBytes, &newTarget)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEqual(t, time.Time{}, newTarget.CreatedAt)
}

func TestTargetPostInvalidTarget(t *testing.T) {
	r := Serve(&testutil.Mock{})
	w := httptest.NewRecorder()
	jsonBody := "{}"
	req, _ := http.NewRequest(http.MethodPost, "/targets",
		strings.NewReader(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTargetPostDbError(t *testing.T) {
	dbMock := testutil.Mock{
		InsertTargetError: errors.New("Connection error"),
	}
	r := Serve(&dbMock)
	w := httptest.NewRecorder()
	target := getValidTarget()
	jsonBody, _ := json.Marshal(target)
	req, _ := http.NewRequest(http.MethodPost, "/targets",
		strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTargetPostErrorDuplicate(t *testing.T) {
	dbMock := testutil.Mock{
		InsertTargetError: errors.New("Duplicate"),
	}
	r := Serve(&dbMock)
	w := httptest.NewRecorder()
	target := getValidTarget()
	jsonBody, _ := json.Marshal(target)
	req, _ := http.NewRequest(http.MethodPost, "/targets",
		strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	jsonBytes, _ := ioutil.ReadAll(w.Result().Body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Duplicate target ID")
}

func TestTargetPostErrorInvalidID(t *testing.T) {
	dbMock := testutil.Mock{}
	r := Serve(&dbMock)
	w := httptest.NewRecorder()
	target := getValidTarget()
	target.ID = ".target/"
	jsonBody, _ := json.Marshal(target)
	req, _ := http.NewRequest(http.MethodPost, "/targets",
		strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	jsonBytes, _ := ioutil.ReadAll(w.Result().Body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Invalid target ID")
}

func TestTargetPostNoParentFoundError(t *testing.T) {
	r := Serve(&testutil.Mock{})
	w := httptest.NewRecorder()
	target := getValidTarget()
	target.ID = "mytargets.target-1"
	jsonBody, _ := json.Marshal(target)
	req, _ := http.NewRequest(http.MethodPost, "/targets",
		strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	jsonBytes, _ := ioutil.ReadAll(w.Result().Body)
	fmt.Println(string(jsonBytes))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, string(jsonBytes), "Target parent does not exist")
}

func TestTargetPostSearchingForParentError(t *testing.T) {
	r := Serve(&testutil.Mock{
		GetTargetByIDError: errors.New("Connection error"),
	})
	w := httptest.NewRecorder()
	target := getValidTarget()
	target.ID = "mytargets.target-1"
	jsonBody, _ := json.Marshal(target)
	req, _ := http.NewRequest(http.MethodPost, "/targets",
		strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTargetPostParentExists(t *testing.T) {
	r := Serve(&testutil.Mock{
		GetTargetByIDResult: &types.Target{
			ID: "mytargets",
		},
	})
	w := httptest.NewRecorder()
	target := getValidTarget()
	target.ID = "mytargets.target-1"
	jsonBody, _ := json.Marshal(target)
	req, _ := http.NewRequest(http.MethodPost, "/targets",
		strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	jsonBytes, _ := ioutil.ReadAll(w.Result().Body)
	var newTarget types.Target
	json.Unmarshal(jsonBytes, &newTarget)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEqual(t, time.Time{}, newTarget.CreatedAt)
}

func getValidTarget() *types.Target {
	return &types.Target{
		ID:          "mytarget-1",
		Description: "my target",
	}
}
