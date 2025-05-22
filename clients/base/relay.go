package base

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Relay(ctx context.Context, method, targetURL string, body io.ReadCloser, header http.Header, client *http.Client) (io.ReadCloser, http.Header, error) {
	resp, err := RelayHttpRequest(ctx, method, targetURL, body, header, client)
	if err != nil {
		return nil, nil, err
	}
	return resp.Body, resp.Header, nil
}

func RelayHttpRequest(ctx context.Context, method, targetURL string, body io.ReadCloser, header http.Header, client *http.Client) (*http.Response, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("parse url error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("new request error: %w", err)
	}
	req.Header = header
	return RelayRequest(req, client)
}

func RelayRequest(req *http.Request, client *http.Client) (*http.Response, error) {
	return client.Do(req)
}
