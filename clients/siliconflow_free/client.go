package siliconflow_free

import (
	"context"
	"github.com/jiu-u/oai-adapter/clients/openai"
	"net/url"
)

type Client struct {
	*openai.Client
}

func NewClient(EndPoint, apiKey string, proxy *url.URL) *Client {
	return &Client{
		Client: openai.NewClient(EndPoint, apiKey, proxy),
	}
}

func (c *Client) Models(ctx context.Context) ([]string, error) {
	return FreeModels, nil
}
