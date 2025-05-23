package oai_adapter

import (
	"context"
	v1 "github.com/jiu-u/oai-adapter/api/v1"
	"github.com/jiu-u/oai-adapter/clients/base"
	"io"
	"net/http"
)

var _ Adapter = (*base.Client)(nil)

type Adapter interface {
	// Settings
	SetClient(client *http.Client)
	// Relay
	RelayRequest(ctx context.Context, method, targetPath string, body io.Reader, header http.Header) (io.ReadCloser, http.Header, error)
	// Responses
	CreateResponses(ctx context.Context, req *v1.ResponsesRequest) (io.ReadCloser, http.Header, error)
	// ChatCompletions
	CreateChatCompletions(ctx context.Context, req *v1.ChatCompletionRequest) (io.ReadCloser, http.Header, error)
	// Completions(Legacy)
	CreateCompletions(ctx context.Context, req *v1.CompletionsRequest) (io.ReadCloser, http.Header, error)
	// Embeddings
	CreateEmbeddings(ctx context.Context, req *v1.EmbeddingsRequest) (io.ReadCloser, http.Header, error)
	// rerank
	CreateRerank(ctx context.Context, req *v1.RerankRequest) (io.ReadCloser, http.Header, error)
	// Image
	CreateImage(ctx context.Context, req *v1.ImageGenerateRequest) (io.ReadCloser, http.Header, error)
	CreateImageEdit(ctx context.Context, req *v1.ImageEditRequest) (io.ReadCloser, http.Header, error)
	CreateImageVariation(ctx context.Context, req *v1.ImageVariationRequest) (io.ReadCloser, http.Header, error)
	// Audio
	CreateSpeech(ctx context.Context, req *v1.AudioSpeechRequest) (io.ReadCloser, http.Header, error)
	CreateTranslation(ctx context.Context, req *v1.TranslationRequest) (io.ReadCloser, http.Header, error)
	CreateTranscription(ctx context.Context, req *v1.TranscriptionRequest) (io.ReadCloser, http.Header, error)
	// Video
	CreateVideoSubmit(ctx context.Context, req *v1.VideoRequest) (*v1.VideoResponse, error)
	GetVideoStatus(ctx context.Context, externalID string) (bool, any, error)
	// Models
	Models(ctx context.Context) (*v1.ModelResponse, error)
	// todo
	// realtime
}
