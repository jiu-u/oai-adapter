package ollama_oai

import (
	"context"
	v1 "github.com/jiu-u/oai-adapter/api/v1"
	"github.com/jiu-u/oai-adapter/clients/base"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	*base.Client
}

func NewClient(endPoint, apiKey string) *Client {
	endPoint = strings.TrimSpace(endPoint)
	endPoint = strings.TrimRight(endPoint, "/")
	endPoint = endPoint + "/v1"
	return &Client{
		Client: base.NewClient(endPoint, apiKey),
	}
}

func (c *Client) CreateRerank(ctx context.Context, req *v1.RerankRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod()
}

func (c *Client) CreateSpeech(ctx context.Context, req *v1.AudioSpeechRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod()
}

func (c *Client) CreateTranslation(ctx context.Context, req *v1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod()
}

func (c *Client) CreateTranscription(ctx context.Context, req *v1.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod()
}

func (c *Client) CreateImage(ctx context.Context, req *v1.ImageGenerateRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod()
}

func (c *Client) CreateImageEdit(ctx context.Context, req *v1.ImageEditRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod()
}

func (c *Client) CreateImageVariation(ctx context.Context, req *v1.ImageVariationRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod()
}

func (c *Client) CreateVideoSubmit(ctx context.Context, req *v1.VideoRequest) (*v1.VideoResponse, error) {
	return nil, v1.NoImplementError
}

func (c *Client) GetVideoStatus(ctx context.Context, externalID string) (bool, any, error) {
	return false, nil, v1.NoImplementError
}
