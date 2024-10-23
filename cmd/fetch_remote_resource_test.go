package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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

		data, err := fetchRemoteResource(testConfig)
		if err != nil && err.Error() != tc.err.Error() {
			t.Fatal(err)
		}

		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}
	}
}
