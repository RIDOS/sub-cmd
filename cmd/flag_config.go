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
	url      string
	verb     string
	body     []byte
	formData []string
	filePath string
	upload   string
	Bytes    io.Reader
}

type arrayFlags []string

func (a *arrayFlags) String() string {
	return fmt.Sprintf("%v", *a)
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
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

func flagConfig(w io.Writer, args []string) (httpConfig, error) {
	hc := httpConfig{}

	var (
		httpVerb     string
		filePath     string
		bodyFlag     string
		bodyFileFlag string
		formDataFlag arrayFlags
		uploadFlag   string
	)

	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&httpVerb, "verb", "GET", "HTTP method")
	fs.StringVar(&filePath, "o", "", "Wtite response in file "+outputFileName)
	fs.StringVar(&bodyFlag, "body", "", "Write body form-data for request (format: json)")
	fs.StringVar(&bodyFileFlag, "body-file", "", "File path for request (format file: json)")
	fs.Var(&formDataFlag, "form-data", "Form data params (format: name=value)")
	fs.StringVar(&uploadFlag, "upload", "", "The path to the file to send files using the POST method")

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
		return hc, err
	}

	if fs.NArg() != 1 {
		return hc, ErrNoServerSpecified
	}

	if !isValidHttpMethod(httpVerb) {
		return hc, ErrInvalidMethod
	}

	if len(bodyFlag) > 0 {
		data, err := parseAndReformatJson(bodyFlag)
		if err != nil {
			return hc, ErrInvalidEncoding
		}
		hc.body = []byte(data)
	}
	if len(bodyFileFlag) > 0 {
		file, err := os.Open(bodyFileFlag)
		if err != nil {
			return hc, err
		}
		defer file.Close()
		if hc.body, err = io.ReadAll(file); err != nil {
			return hc, err
		}
	}

	if len(uploadFlag) > 0 {
		fileData, err := os.Open(uploadFlag)
		if err != nil {
			return hc, err
		}
		defer fileData.Close()
		hc.upload = uploadFlag
		hc.Bytes = fileData
	}

	hc.url = fs.Arg(0)
	hc.verb = httpVerb
	hc.filePath = filePath
	hc.formData = formDataFlag

	return hc, nil
}
