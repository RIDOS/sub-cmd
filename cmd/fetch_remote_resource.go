package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
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
		contentType, body, err := hc.preparePostData()
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(body)
		r, err = http.Post(hc.url, contentType, reader)
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

func (hc *httpConfig) preparePostData() (string, []byte, error) {
	var dct string = "application/json"
	if len(hc.body) > 0 {
		return dct, hc.body, nil
	}

	if len(hc.formData) > 0 {
		var b bytes.Buffer
		var err error
		var fw io.Writer

		mw := multipart.NewWriter(&b)

		for _, data := range hc.formData {
			splitedData := strings.Split(data, "=")
			if len(splitedData) < 2 {
				return "", []byte{}, fmt.Errorf("invalid form data: %s", data)
			}
			fw, err = mw.CreateFormField(splitedData[0])
			if err != nil {
				return "", []byte{}, err
			}
			fmt.Fprintf(fw, splitedData[1])
		}

		contentType := mw.FormDataContentType()

		return contentType, b.Bytes(), nil
	}

	return "", []byte{}, errors.New("Prepare post data fale: Config is empty")
}
