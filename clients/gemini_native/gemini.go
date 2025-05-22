package gemini_native

import (
	"github.com/jiu-u/oai-adapter/clients/base"
	"github.com/jiu-u/oai-adapter/constant"
	"strings"
)

var GeminiVersion = "v1beta"

type Client struct {
	*base.Client
}

func NewClient(endPoint, apiKey string) *Client {
	if endPoint == "" {
		endPoint = constant.GeminiDefaultURL
	}
	endPoint = strings.TrimSpace(endPoint)
	endPoint = strings.TrimRight(endPoint, "/")
	endPoint = endPoint + "/" + GeminiVersion
	return &Client{
		Client: base.NewClient(endPoint, apiKey),
	}
}
