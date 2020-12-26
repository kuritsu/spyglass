package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuritsu/spyglass/api/storage/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTargetGet(t *testing.T) {
	r := Serve(&testutil.Mock{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/targets", nil)
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
