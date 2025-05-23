package tools

import (
	"io"
	"net/http"
)

func GetReadCloserFromURL(downloadURL string) (io.ReadCloser, http.Header, error) {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != 200 {
		return nil, nil, err
	}
	header := http.Header{}
	if resp.Header.Get("Content-Type") != "" {
		header.Set("Content-Type", resp.Header.Get("Content-Type"))
	} else {
		header.Set("Content-Type", "application/octet-stream")
	}
	return resp.Body, header, nil
}
