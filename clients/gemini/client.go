package gemini

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	stdurl "net/url"
	"time"
)

type Client struct {
	EndPoint string
	APIKey   string
	ProxyURL *stdurl.URL
}

func NewClient(EndPoint, apiKey string, proxy *stdurl.URL) *Client {
	return &Client{
		EndPoint: EndPoint,
		APIKey:   apiKey,
		ProxyURL: proxy,
	}
}

type ReqPoint = string

const (
	ReqPointChat       ReqPoint = "chat"
	ReqPointChatStream ReqPoint = "stream"
	ReqPointModel      ReqPoint = "model"
	ReqPointEmbedding  ReqPoint = "embedding"
)

func (c *Client) DoRequest(ctx context.Context, url string, Method string, body io.Reader, contextType string) (io.ReadCloser, http.Header, error) {
	request, err := http.NewRequestWithContext(ctx, Method, url, body)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Set("Content-Type", contextType)
	//request.Header.Set("Authorization", "Bearer "+c.APIKey)
	client := &http.Client{Timeout: 30 * time.Second}
	if c.ProxyURL != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(c.ProxyURL),
		}
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		errDetail := ""
		errBytes, err := io.ReadAll(resp.Body)
		if err == nil {
			errDetail = string(errBytes)
		}
		return io.NopCloser(bytes.NewReader([]byte(errDetail))), nil, errors.New(resp.Status + ": " + errDetail)
	}
	return resp.Body, resp.Header, nil
}

func (c *Client) DoJsonRequest(ctx context.Context, Method string, url string, body io.Reader) (io.ReadCloser, http.Header, error) {
	return c.DoRequest(ctx, url, Method, body, "application/json")
}

func (c *Client) DoFormRequest(ctx context.Context, url string, body io.Reader, contentType string) (io.ReadCloser, http.Header, error) {
	// 创建 multipart 写入器
	return c.DoRequest(ctx, url, "POST", body, contentType)
}

func (c *Client) DoOAIRequest(ctx context.Context, url string, Method string, body io.Reader, contextType string) (io.ReadCloser, http.Header, error) {
	request, err := http.NewRequestWithContext(ctx, Method, url, body)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Set("Content-Type", contextType)
	request.Header.Set("Authorization", "Bearer "+c.APIKey)
	client := &http.Client{Timeout: 30 * time.Second}
	if c.ProxyURL != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(c.ProxyURL),
		}
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		errDetail := ""
		errBytes, err := io.ReadAll(resp.Body)
		if err == nil {
			errDetail = string(errBytes)
		}
		return io.NopCloser(bytes.NewReader([]byte(errDetail))), nil, errors.New(resp.Status + ": " + errDetail)
	}
	return resp.Body, resp.Header, nil
}

func (c *Client) DoOAIJsonRequest(ctx context.Context, Method string, url string, body io.Reader) (io.ReadCloser, http.Header, error) {
	return c.DoRequest(ctx, url, Method, body, "application/json")
}

func (c *Client) DoOAIFormRequest(ctx context.Context, url string, body io.Reader, contentType string) (io.ReadCloser, http.Header, error) {
	// 创建 multipart 写入器
	return c.DoRequest(ctx, url, "POST", body, contentType)
}
