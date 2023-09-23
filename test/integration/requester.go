package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type requester struct {
	url     string
	payload interface{}
	method  string
	headers map[string]string
}

func getReader(payload interface{}) (io.Reader, error) {
	bytesPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bytesPayload), nil
}

func (r requester) do() (*http.Response, error) {
	var (
		reader io.Reader
		err    error
	)

	if r.payload != nil {
		reader, err = getReader(r.payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(r.method, r.url, reader)
	if err != nil {
		return nil, err
	}

	for k, v := range r.headers {
		req.Header.Add(k, v)
	}

	return http.DefaultClient.Do(req)
}
