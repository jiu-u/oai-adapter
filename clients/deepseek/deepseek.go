package deepseek

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
		endPoint = constant.DeepSeekDefaultURL
	}
	endPoint = strings.TrimSpace(endPoint)
	endPoint = strings.TrimRight(endPoint, "/")
	client := base2.NewClient(endPoint, apiKey)
	return &Client{
		Client: client,
	}
}

func (c *Client) CreateChatCompletions(ctx context.Context, req *v1.ChatCompletionRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/v1/chat/completions"
	return c.SamePostJob(ctx, targetUrl, req, "application/json")
}

func (c *Client) CreateCompletions(ctx context.Context, req *v1.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/beta/completions"
	return c.SamePostJob(ctx, targetUrl, req, "application/json")
}

func (c *Client) CreateResponses(ctx context.Context, req *v1.ResponsesRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod()
}

func (c *Client) CreateEmbeddings(ctx context.Context, req *v1.EmbeddingsRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod()
}

func (c *Client) CreateRerank(ctx context.Context, req *v1.RerankRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod()
}

func (c *Client) CreateImage(ctx context.Context, req *v1.ImageGenerateRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod()
}

func (c *Client) CreateImageEdit(ctx context.Context, req *v1.ImageEditRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod()
}

func (c *Client) CreateImageVariation(ctx context.Context, req *v1.ImageVariationRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod()
}

func (c *Client) CreateSpeech(ctx context.Context, req *v1.AudioSpeechRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod()
}

func (c *Client) CreateTranslation(ctx context.Context, req *v1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod()
}

func (c *Client) CreateTranscription(ctx context.Context, req *v1.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	return base2.NoImplementMethod()
}

func (c *Client) CreateVideoSubmit(ctx context.Context, req *v1.VideoRequest) (*v1.VideoResponse, error) {
	return nil, v1.NoImplementError
}

func (c *Client) GetVideoStatus(ctx context.Context, externalID string) (bool, any, error) {
	return false, nil, v1.NoImplementError
}
