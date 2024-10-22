package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func fetchRemoteResource(hc httpConfig) ([]byte, error) {
	var err error
	var r *http.Response

	switch hc.verb {
	case http.MethodGet:
		r, err = http.Get(hc.url)
	case http.MethodHead:
		r, err = http.Head(hc.url)
	case http.MethodPost:
		reader := bytes.NewReader(hc.body)
		r, err = http.Post(hc.url, "application/json", reader)
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
