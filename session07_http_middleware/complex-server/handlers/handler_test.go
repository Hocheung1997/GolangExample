package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hocheung1997/complex-server/config"
)

func TestApiHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/api", nil)
	w := httptest.NewRecorder()
	b := new(bytes.Buffer)
	c := config.InitConfig(b)

	apiHandler(w, r, c)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status: %v, Got: %v\n", http.StatusOK, err)
	}

	expectedResponseBody := "Hello world!"

	if expectedResponseBody != string(body) {
		t.Errorf("Expected response body string: %s, Got: %s\n", expectedResponseBody, string(body))
	}

}
