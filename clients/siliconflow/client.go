package siliconflow

import (
	"github.com/jiu-u/oai-adapter/clients/openai"
	stdurl "net/url"
)

type Client struct {
	*openai.Client
}

func NewClient(EndPoint, apiKey string, proxy *stdurl.URL) *Client {
	return &Client{
		Client: openai.NewClient(EndPoint, apiKey, proxy),
	}
}
