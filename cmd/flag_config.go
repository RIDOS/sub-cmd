package cmd

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
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

func (hc *httpConfig) preparePostData() (string, []byte, error) {
	var cxt string = "application/json"
	if len(hc.body) > 0 {
		return cxt, hc.body, nil
	}

	if len(hc.formData) > 0 || len(hc.upload) > 0 {
		var b bytes.Buffer
		var err error
		var fw io.Writer

		mw := multipart.NewWriter(&b)

		if len(hc.formData) > 0 {
			for _, data := range hc.formData {
				splitedData, err := strToParams(data)
				if err != nil {
					return "", nil, err
				}
				fw, err = mw.CreateFormField(splitedData[0])
				if err != nil {
					return "", nil, err
				}
				fmt.Fprintf(fw, splitedData[1])
			}
		}

		if len(hc.upload) > 0 {
			fw, err = mw.CreateFormFile("filedata", hc.upload)
			if err != nil {
				return "", nil, err
			}

			_, err = io.Copy(fw, hc.Bytes)
			if err != nil {
				return "", nil, err
			}
		}

		err = mw.Close()
		if err != nil {
			return "", nil, err
		}

		contentType := mw.FormDataContentType()

		return contentType, b.Bytes(), nil
	}

	return "", []byte{}, errors.New("Prepare post data fale: Config is empty")
}
