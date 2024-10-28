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
			args: []string{"-verb", "DELETE", "http://localhost"},
			err:  ErrInvalidMethod,
		},
		{
			args: []string{"-verb", "PUT", "http://localhost"},
			err:  ErrInvalidMethod,
		},
		{
			args: []string{"-verb", "UPDATE", "http://localhost"},
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
