package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecureHeaders(t *testing.T) {
	t.Parallel() // run this concurrently just an experiment

	request := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close() // important, otherwise the TCP connection cannot be reused after body is not used
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "deny", response.Header.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", response.Header.Get("X-XSS-Protection"))
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "OK", string(body))
}
