package oai_no_models

import (
	"context"
	"github.com/jiu-u/oai-adapter/clients/openai"
	stdurl "net/url"
)

type Client struct {
	*openai.Client
	models []string
}

func NewClient(apiKey string, EndPoint string, proxy *stdurl.URL, models []string) *Client {
	return &Client{
		Client: openai.NewClient(EndPoint, apiKey, proxy),
		models: models,
	}
}

func (c *Client) Models(ctx context.Context) ([]string, error) {
	return c.models, nil
}
