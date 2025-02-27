package deepseek

import (
	"bytes"
	"context"
	"github.com/bytedance/sonic"
	"github.com/jiu-u/oai-adapter/api"
	"github.com/jiu-u/oai-adapter/clients/openai"
	"github.com/jiu-u/oai-adapter/constant"
	"io"
	"net/http"
	stdurl "net/url"
)

type Client struct {
	*openai.Client
}

func NewClient(endpoint, apiKey string, proxy *stdurl.URL) *Client {
	if len(endpoint) == 0 {
		endpoint = constant.DeepSeekDefaultURL
	}
	return &Client{
		Client: openai.NewClient(endpoint, apiKey, proxy),
	}
}

func (c *Client) Completions(ctx context.Context, req *api.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	url := c.EndPoint + "/beta/completions"
	return c.DoJsonRequest(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
}

func (c *Client) CompletionsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	url := c.EndPoint + "/beta/completions"
	return c.DoJsonRequest(ctx, "POST", url, bytes.NewBuffer(req))
}

func (c *Client) Embeddings(ctx context.Context, req *api.EmbeddingRequest) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}

func (c *Client) EmbeddingsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}

func (c *Client) CreateSpeech(ctx context.Context, req *api.SpeechRequest) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}

func (c *Client) CreateSpeechByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}

func (c *Client) Transcriptions(ctx context.Context, req *api.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}

func (c *Client) Translations(ctx context.Context, req *api.TranslationRequest) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}

func (c *Client) CreateImage(ctx context.Context, req *api.CreateImageRequest) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}

func (c *Client) CreateImageByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}

func (c *Client) CreateImageEdit(ctx context.Context, req *api.EditImageRequest) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}

func (c *Client) ImageVariations(ctx context.Context, req *api.CreateImageVariationRequest) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}
