package cmd

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	MethodGet  = "GET"
	MethodPost = "POST"
	MethodHead = "HEAD"
)

var outputFileName string = "output.html"
var httpMethods = []string{MethodGet, MethodPost, MethodHead}

type httpConfig struct {
	url  string
	verb string
}

type Output interface {
	Write(data []byte) error
}

type ConsoleOutput struct {
	wirter io.Writer
}

func (c *ConsoleOutput) Write(data []byte) error {
	_, err := fmt.Fprintf(c.wirter, "%s\n", data)
	if err != nil {
		return err
	}
	return nil
}

type FileOutput struct {
	filePath string
}

func (fileOutput *FileOutput) Write(data []byte) error {
	err := os.WriteFile(fileOutput.filePath+outputFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func HandleHttp(w io.Writer, args []string) error {
	var httpVerb string
	var filePath string
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&httpVerb, "verb", "GET", "HTTP method")
	fs.StringVar(&filePath, "o", "", "Wtite response in file "+outputFileName)

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

func fetchRemoteResource(hc httpConfig) ([]byte, error) {
	var err error
	var r *http.Response

	switch hc.verb {
	case MethodGet:
		r, err = http.Get(hc.url)
	case MethodHead:
		r, err = http.Head(hc.url)
	case MethodPost:
		r, err = http.Post(hc.url, "application/json", nil)
	default:
		err = ErrInvalidMethod
	}

	if err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unxpected status code: %d", r.StatusCode)
	}

	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
