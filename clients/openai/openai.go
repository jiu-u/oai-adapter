package openai

import (
	"context"
	v1 "github.com/jiu-u/oai-adapter/api/v1"
	base2 "github.com/jiu-u/oai-adapter/clients/base"
	"github.com/jiu-u/oai-adapter/constant"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	*base2.Client
}

func NewClient(endPoint, apiKey string) *Client {
	if endPoint == "" {
		endPoint = constant.OpenAIDefaultURL
	}
	endPoint = strings.TrimSpace(endPoint)
	endPoint = strings.TrimRight(endPoint, "/")
	endPoint = endPoint + "/v1"
	client := base2.NewClient(endPoint, apiKey)
	return &Client{
		Client: client,
	}
}

func (c *Client) CreateRerank(ctx context.Context, req *v1.RerankRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod(ctx, req)
}

func (c *Client) CreateVideoSubmit(ctx context.Context, req *v1.VideoRequest) (*v1.VideoResponse, error) {
	_, _, err := base2.NoImplementMethod(ctx, req)
	return nil, err
}

func (c *Client) GetVideoStatus(ctx context.Context, externalID string) (bool, any, error) {
	_, _, err := base2.NoImplementMethod(ctx, externalID)
	return false, nil, err
}
