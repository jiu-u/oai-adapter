package oai_adapter

import (
	"context"
	"github.com/jiu-u/oai-adapter/api"
	"io"
	"net/http"
)

type Adapter interface {
	ChatCompletions(ctx context.Context, req *api.ChatRequest) (io.ReadCloser, http.Header, error)
	ChatCompletionsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error)
	Models(ctx context.Context) ([]string, error)
	Completions(ctx context.Context, req *api.CompletionsRequest) (io.ReadCloser, http.Header, error)
	CompletionsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error)
	Embeddings(ctx context.Context, req *api.EmbeddingRequest) (io.ReadCloser, http.Header, error)
	EmbeddingsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error)
	CreateSpeech(ctx context.Context, req *api.SpeechRequest) (io.ReadCloser, http.Header, error)
	CreateSpeechByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error)
	Transcriptions(ctx context.Context, req *api.TranscriptionRequest) (io.ReadCloser, http.Header, error)
	Translations(ctx context.Context, req *api.TranslationRequest) (io.ReadCloser, http.Header, error)
	CreateImage(ctx context.Context, req *api.CreateImageRequest) (io.ReadCloser, http.Header, error)
	CreateImageByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error)
	CreateImageEdit(ctx context.Context, req *api.EditImageRequest) (io.ReadCloser, http.Header, error)
	ImageVariations(ctx context.Context, req *api.CreateImageVariationRequest) (io.ReadCloser, http.Header, error)
}
