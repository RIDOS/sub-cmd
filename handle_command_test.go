package main

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"testing"
	"time"

	"github.com/RIDOS/sub-cmd/cmd"
)

func TestHanldeCommand(t *testing.T) {
	usgaeMessage := `Usage: mync [http|grpc] -h

http: A HTTP client.
http: <options> server

Options: 
  -body string
    	Write body form-data for request (format: json)
  -body-file string
    	File path for request (format file: json)
  -disable-redirect
    	Disable redirect for response
  -form-data value
    	Form data params (format: name=value)
  -header value
    	Request Headers (format: name=value)
  -o string
    	Wtite response in file output.html
  -upload string
    	The path to the file to send files using the POST method
  -verb string
    	HTTP method (default "GET")

grpc: A gRPC client.
grpc: <options> server

Options: 
  -body string
    	Body of request
  -method string
    	Method to call
`

	testConfig := []struct {
		args     []string
		output   string
		err      error
		exitCode int
	}{
		// Test when application start whitout argunets.
		{
			args:     []string{},
			err:      errInvalidSubCommand,
			output:   "Invalid sub-command specified\n" + usgaeMessage,
			exitCode: 1,
		},
		// Test when use "-h" in arguments.
		{
			args:     []string{"-h"},
			err:      nil,
			output:   usgaeMessage,
			exitCode: 0,
		},
		// Test when use something in arguments.
		{
			args:     []string{"foo"},
			err:      errInvalidSubCommand,
			output:   "Invalid sub-command specified\n" + usgaeMessage,
			exitCode: 1,
		},
		// Test http app.
		{
			args:     []string{"http", "--help"},
			err:      errors.New("flag: help requested"),
			exitCode: 1,
		},
		// Test grpc app.
		{
			args:     []string{"grpc", "foo"},
			err:      nil,
			output:   "Execution grpc command\n",
			exitCode: 0,
		},
		{
			args:     []string{"http", "-verb", "DELETE", "http://localhost"},
			err:      cmd.ErrInvalidMethod,
			exitCode: 1,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfig {
		err := handleCommand(ctx, byteBuf, tc.args)

		if tc.err == nil && err != nil {
			t.Fatalf("Expected nill error, got \n%v\n", err)
		}

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error \n%v\n,got \n%v\n", tc.err, err)
		}

		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Errorf("Expected output to be: \n%#v\n, Got: \n%#v\n", tc.output, gotOutput)
			}
		}
		byteBuf.Reset()

		// Test exit code.
		cmd := exec.Command("./console-prog", tc.args...)
		err = cmd.Run()

		if err != nil {
			var exitErr *exec.ExitError
			// This is for GitHub Actions.
			if errors.As(err, &exitErr) {
				gotExitCode := exitErr.ExitCode()
				if tc.exitCode != gotExitCode {
					t.Fatalf("Expected exit code to be: %v, Got %v\n", tc.exitCode, gotExitCode)
				}
			}
		}
	}
}
