package testutil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

// MakeRequest to an http endpoint
func MakeRequest(method string, url string, body interface{}, r *gin.Engine) (*httptest.ResponseRecorder, []byte) {
	w := httptest.NewRecorder()
	var bodyReader io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = strings.NewReader(string(jsonBody))
	}

	req, _ := http.NewRequest(method, url, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	jsonBytes, _ := ioutil.ReadAll(w.Result().Body)
	return w, jsonBytes
}
