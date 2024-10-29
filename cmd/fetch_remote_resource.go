package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func fetchRemoteResource(ctx context.Context, client *http.Client, hc httpConfig) ([]byte, error) {
	var (
		err error
		req *http.Request
	)

	switch hc.verb {
	case http.MethodGet:
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, hc.url, nil)
	case http.MethodHead:
		req, err = http.NewRequestWithContext(ctx, http.MethodHead, hc.url, nil)
	case http.MethodPost:
		contentType, body, err := hc.preparePostData()
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(body)
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, hc.url, reader)
		req.Header.Set("Content-Type", contentType)
	default:
		err = ErrInvalidMethod
	}

	if err != nil {
		return nil, err
	}

	if len(hc.headers) > 0 {
		for k, v := range hc.headers {
			req.Header.Add(k, v)
		}
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unxpected status code: %d", r.StatusCode)
	}

	return io.ReadAll(r.Body)
}
