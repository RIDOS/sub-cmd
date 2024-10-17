package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startTestHTTPServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello World!")
		}))
	return ts
}

func TestFetchRemoteResource(t *testing.T) {
	ts := startTestHTTPServer()
	defer ts.Close()

	expected := "Hello World!"

	testConfig := httpConfig{url: ts.URL, verb: "GET"}

	data, err := fetchRemoteResource(testConfig)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}
