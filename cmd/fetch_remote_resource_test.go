package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func startTestHTTPServer(response string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodHead {
				w.Header().Set("Content-Length", fmt.Sprint(len(response)))
			}
			fmt.Fprint(w, response)
		}))
	return ts
}

func TestFetchRemoteResource(t *testing.T) {
	testConfig := []struct {
		request     string
		response    string
		method      string
		err         error
		contentType string
	}{
		{
			request:  "Hello World!",
			response: "Hello World!",
			method:   "GET",
		},
		// This method not respone ther head.
		{
			request:  "",
			response: "Hello World!",
			method:   "HEAD",
		},
		{
			request:     "Hello World!",
			response:    "Hello World!",
			method:      "POST",
			contentType: "body",
		},
		{
			request:     "nameRichard",
			method:      "POST",
			contentType: "formData",
			err:         errors.New("invalid form data: nameRichard"),
		},
		{
			request:     "name=Richard",
			response:    "name=Richard",
			method:      "POST",
			contentType: "formData",
		},
		{
			request:  "",
			response: "Hello World!",
			method:   "DELETE",
			err:      ErrInvalidMethod,
		},
		{
			request:  "",
			response: "Hello World!",
			method:   "UPDATE",
			err:      ErrInvalidMethod,
		},
		{
			request:  "",
			response: "Hello World!",
			method:   "PUT",
			err:      ErrInvalidMethod,
		},
	}

	for _, tc := range testConfig {
		ts := startTestHTTPServer(tc.response)
		defer ts.Close()

		expected := tc.request

		testConfig := httpConfig{url: ts.URL, verb: tc.method}
		if tc.contentType == "body" {
			testConfig.body = []byte(tc.request)
		}
		if tc.contentType == "formData" {
			testConfig.formData = []string{tc.request}
		}

		client := createHTTPClientWithTimeout(withTimeOut(20 * time.Millisecond))
		data, err := fetchRemoteResource(context.Background(), client, testConfig)
		if err != nil && err.Error() != tc.err.Error() {
			t.Errorf("Expected %s, got %s", tc.err, err)
		}

		if string(data) != expected && len(tc.err.Error()) == 0 {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}
	}
}

func startBadTestHTTPServerV2(shutdownServer chan struct{}) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-shutdownServer
		fmt.Fprint(w, "Hello World")
	}))
	return ts
}

func TestFetchBadRemoteResourceV2(t *testing.T) {
	shutdownServer := make(chan struct{})
	ts := startBadTestHTTPServerV2(shutdownServer)
	defer ts.Close()

	defer func() {
		shutdownServer <- struct{}{}
	}()

	client := createHTTPClientWithTimeout(withTimeOut(20 * time.Millisecond))
	hc := httpConfig{
		url:  ts.URL,
		verb: "GET",
	}
	_, err := fetchRemoteResource(context.Background(), client, hc)
	if err == nil {
		t.Fatalf("Expected not-nil error")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("Expected error to contain: context deadline exceeded, Got: %v", err.Error())
	}
}
