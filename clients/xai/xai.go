package xai

import (
	"github.com/jiu-u/oai-adapter/clients/openai"
	"github.com/jiu-u/oai-adapter/constant"
	stdurl "net/url"
)

type Client struct {
	*openai.Client
}

func NewClient(EndPoint, apiKey string, proxy *stdurl.URL) *Client {
	if EndPoint == "" {
		EndPoint = constant.XAIDefaultURL
	}
	return &Client{
		Client: openai.NewClient(EndPoint, apiKey, proxy),
	}
}
