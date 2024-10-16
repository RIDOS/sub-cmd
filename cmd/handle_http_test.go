package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestHandleHttp(t *testing.T) {
	usageMessage := `
http: A HTTP client.
http: <options> server

Options: 
  -verb string
    	HTTP method (default "GET")
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
			args:   []string{"http://localhost"},
			err:    nil,
			output: "Execution http command\n",
		},
		{
			args:   []string{"-verb", "GET", "http://localhost"},
			err:    nil,
			output: "Execution http command\n",
		},
		{
			args:   []string{"-verb", "POST", "http://localhost"},
			err:    nil,
			output: "Execution http command\n",
		},
		{
			args:   []string{"-verb", "HEAD", "http://localhost"},
			err:    nil,
			output: "Execution http command\n",
		},
		{
			args: []string{"-verb", "DELETE", "http://localhost"},
			err:  ErrInvalidMethod,
		},
		{
			args:   []string{"-verb", "PUT", "http://localhost"},
			err:  ErrInvalidMethod,
		},
		{
			args:   []string{"-verb", "UPDATE", "http://localhost"},
			err:  ErrInvalidMethod,
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfig {
		err := HandleHttp(byteBuf, tc.args)

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
