package gemini_oai

import (
	"context"
	v1 "github.com/jiu-u/oai-adapter/api/v1"
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
	endPoint = endPoint + "/" + GeminiVersion + "/openai"
	return &Client{
		Client: base.NewClient(endPoint, apiKey),
	}
}

func NewClientWithVersion(endPoint, apiKey, version string) *Client {
	if endPoint == "" {
		endPoint = constant.GeminiDefaultURL
	}
	endPoint = strings.TrimSpace(endPoint)
	endPoint = strings.TrimRight(endPoint, "/")
	endPoint = endPoint + "/" + version + "/openai"
	return &Client{
		Client: base.NewClient(endPoint, apiKey),
	}
}

func (c *Client) CreateVideoSubmit(ctx context.Context, req *v1.VideoRequest) (*v1.VideoResponse, error) {
	_, _, err := base.NoImplementMethod(ctx, req)
	return nil, err
}

func (c *Client) GetVideoStatus(ctx context.Context, externalID string) (bool, any, error) {
	_, _, err := base.NoImplementMethod(ctx, externalID)
	return false, nil, err
}
