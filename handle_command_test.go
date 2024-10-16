package main

import (
	"bytes"
	"errors"
	"os/exec"
	"testing"

	"github.com/RIDOS/sub-cmd/cmd"
)

func TestHanldeCommand(t *testing.T) {
	usgaeMessage := `Usage: mync [http|grpc] -h

http: A HTTP client.
http: <options> server

Options: 
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
			args:     []string{"http", "foo"},
			err:      nil,
			output:   "Execution http command\n",
			exitCode: 0,
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

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfig {
		err := handleCommand(byteBuf, tc.args)

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
