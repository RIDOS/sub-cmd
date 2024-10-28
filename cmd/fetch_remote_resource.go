package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func fetchRemoteResource(client *http.Client, hc httpConfig) ([]byte, error) {
	var (
		err error
		r   *http.Response
	)

	switch hc.verb {
	case http.MethodGet:
		r, err = client.Get(hc.url)
	case http.MethodHead:
		r, err = client.Head(hc.url)
	case http.MethodPost:
		contentType, body, err := hc.preparePostData()
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(body)
		r, err = client.Post(hc.url, contentType, reader)
	default:
		err = ErrInvalidMethod
	}

	if len(hc.headers) > 0 {
		for k, v := range hc.headers {
			r.Header.Add(k, v)
		}
	}

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unxpected status code: %d", r.StatusCode)
	}

	return io.ReadAll(r.Body)
}
