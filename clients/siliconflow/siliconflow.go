package siliconflow

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
		endPoint = constant.SiliconFlowDefaultURL
	}
	endPoint = strings.TrimSpace(endPoint)
	endPoint = strings.TrimRight(endPoint, "/")
	endPoint = endPoint + "/v1"
	client := base2.NewClient(endPoint, apiKey)
	return &Client{
		Client: client,
	}
}

func (c *Client) CreateCompletions(ctx context.Context, req *v1.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod(ctx, req)
}

func (c *Client) CreateResponses(ctx context.Context, req *v1.ResponsesRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod(ctx, req)
}

func (c *Client) CreateImageEdit(ctx context.Context, req *v1.ImageEditRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod(ctx, req)
}

func (c *Client) CreateImageVariation(ctx context.Context, req *v1.ImageVariationRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod(ctx, req)
}

func (c *Client) CreateTranslation(ctx context.Context, req *v1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod(ctx, req)
}
