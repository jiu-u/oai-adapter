package base

import (
	"context"
	"fmt"
	"github.com/jiu-u/oai-adapter/common"
	"io"
	"net/http"
	"net/url"
)

func Relay(ctx context.Context, method, targetURL string, body io.Reader, header http.Header, client *http.Client) (io.ReadCloser, http.Header, error) {
	resp, err := RelayHttpRequest(ctx, method, targetURL, body, header, client)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode >= 400 {
		return resp.Body, nil, fmt.Errorf("status code error: %d | %w", resp.StatusCode, common.StatusCodeError)
	}
	return resp.Body, resp.Header, nil
}

func RelayHttpRequest(ctx context.Context, method, targetURL string, body io.Reader, header http.Header, client *http.Client) (*http.Response, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("parse url error: %w", err)
	}
	fmt.Println("请求地址：", parsedURL.String())

	req, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("new request error: %w", err)
	}
	for k, v := range header {
		req.Header.Add(k, v[0])
	}
	return RelayRequest(req, client)
}

func RelayRequest(req *http.Request, client *http.Client) (*http.Response, error) {
	return client.Do(req)
}
