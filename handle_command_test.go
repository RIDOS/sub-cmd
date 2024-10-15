package main

import (
	"bytes"
	"testing"
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
		args   []string
		output string
		err    error
	}{
		// Test when application start whitout argunets
		{
			args:   []string{},
			err:    errInvalidSubCommand,
			output: "Invalid sub-command specified\n" + usgaeMessage,
		},
		// Test when use "-h" in arguments
		{
			args:   []string{"-h"},
			err:    nil,
			output: usgaeMessage,
		},
		// Test when use something in arguments
		{
			args:   []string{"foo"},
			err:    errInvalidSubCommand,
			output: "Invalid sub-command specified\n" + usgaeMessage,
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
	}
}
