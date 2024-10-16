package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestHandleGrpc(t *testing.T) {
	usageMessage := `
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
		{
			args: []string{},
			err:  ErrNoServerSpecified,
		},
		{
			args:   []string{"-h"},
			err:    errors.New("flag: help requested"),
			output: usageMessage,
		},
		{
			args:   []string{"--method", "service.host.local/method", "-body", "{}", "https://localhost"},
			err:    nil,
			output: "Execution grpc command\n",
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfig {
		err := HandleGrpc(byteBuf, tc.args)

		// Error test.
		if err != nil {
			if err.Error() != tc.err.Error() {
				t.Errorf("Expected output to be: \n%#v\n, Got: \n%#v\n", tc.err.Error(), err.Error())
			}
		}
		// Output test.
		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Errorf("Expected output to be: \n%#v\n, Got: \n%#v\n", tc.output, gotOutput)
			}
		}
		byteBuf.Reset()
	}
}
