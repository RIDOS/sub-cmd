package cmd

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"slices"
)

var httpMethods = []string{"GET", "POST", "HEAD"}

type httpConfig struct {
	url  string
	verb string
}

func HandleHttp(w io.Writer, args []string) error {
	var v string
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&v, "verb", "GET", "HTTP method")

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

	if !slices.Contains(httpMethods, v) {
		return ErrInvalidMethod
	}

	c := httpConfig{verb: v}
	c.url = fs.Arg(0)

	body, err := fetchRemoteResource(c)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", body)

	return nil
}

func fetchRemoteResource(hc httpConfig) ([]byte, error) {
	var err error
	var r *http.Response

	switch hc.verb {
	case "GET":
		r, err = http.Get(hc.url)
	case "HEAD":
		r, err = http.Head(hc.url)
	case "POST":
		r, err = http.Post(hc.url, "application/json", nil)
	default:
		err = ErrInvalidMethod
	}

	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
