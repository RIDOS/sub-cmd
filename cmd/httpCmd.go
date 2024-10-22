package cmd

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var outputFileName string = "output.html"
var httpMethods = []string{http.MethodGet, http.MethodPost, http.MethodHead}

type httpConfig struct {
	url  string
	verb string
	body []byte
}

func HandleHttp(w io.Writer, args []string) error {
	var (
		httpVerb     string
		filePath     string
		bodyFlag     string
		bodyFileFlag string
	)

	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&httpVerb, "verb", "GET", "HTTP method")
	fs.StringVar(&filePath, "o", "", "Wtite response in file "+outputFileName)
	fs.StringVar(&bodyFlag, "body", "", "Write body form-data for request (format: json)")
	fs.StringVar(&bodyFileFlag, "body-file", "", "File path for request (format file: json)")

	fs.Usage = func() {
		var usageString = `
http: A HTTP client.
http: <options> server`

		fmt.Fprintf(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	if !isValidHttpMethod(httpVerb) {
		return ErrInvalidMethod
	}

	c := httpConfig{verb: httpVerb}
	c.url = fs.Arg(0)

	if len(bodyFlag) > 0 {
		data, err := parseAndReformatJson(bodyFlag)
		if err != nil {
			return ErrInvalidEncoding
		}
		c.body = []byte(data)
	}
	if len(bodyFileFlag) > 0 {
		file, err := os.Open(bodyFileFlag)
		if err != nil {
			return err
		}
		defer file.Close()
		if c.body, err = io.ReadAll(file); err != nil {
			return err
		}
	}

	body, err := fetchRemoteResource(c)
	if err != nil {
		return err
	}

	var output Output
	if filePath != "" {
		output = &FileOutput{filePath: filePath}
	} else {
		output = &ConsoleOutput{wirter: w}
	}

	err = output.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func isValidHttpMethod(method string) bool {
	for _, m := range httpMethods {
		if method == m {
			return true
		}
	}
	return false
}
