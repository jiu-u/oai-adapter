package gemini

import (
	"context"
	"github.com/jiu-u/oai-adapter/api"
	"github.com/jiu-u/oai-adapter/clients/openai"
	"io"
	"net/http"
)

func (c *Client) Completions(ctx context.Context, req *api.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	return openai.NoImplementMethod()
}

func (c *Client) CompletionsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
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
