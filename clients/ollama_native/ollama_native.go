package ollama_native

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	v1 "github.com/jiu-u/oai-adapter/api/v1"
	"github.com/jiu-u/oai-adapter/clients/base"
	"github.com/jiu-u/oai-adapter/tools"
	"github.com/oklog/ulid/v2"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	*base.Client
}

// ollama 能力
// chatCompletions
// completions
// models
// embeddings

func NewClient(endPoint, apiKey string) *Client {
	endPoint = strings.TrimSpace(endPoint)
	endPoint = strings.TrimRight(endPoint, "/")
	return &Client{
		Client: base.NewClient(endPoint, apiKey),
	}
}

func (c *Client) CreateResponses(ctx context.Context, req *v1.ResponsesRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) ConvertChatCompletions(req *v1.ChatCompletionRequest) (io.ReadCloser, error) {
	var err error
	ollamaReq := OllamaChatCompletionRequest{
		Model:    req.Model,
		Messages: []OllamaMessage{},
		Tools:    req.Tools,
		Format:   req.ResponseFormat,
		Options: &Options{
			NumCtx:           int64(req.MaxCompletionTokens),
			Temperature:      req.Temperature,
			Seed:             req.Seed,
			Stop:             nil,
			TopP:             req.TopP,
			PresencePenalty:  req.PresencePenalty,
			FrequencyPenalty: req.FrequencyPenalty,
		},
		Stream:    req.Stream,
		KeepAlive: "",
	}
	for _, msg := range req.Messages {
		if msg.IsStringContent() {
			ollamaReq.Messages = append(ollamaReq.Messages, OllamaMessage{
				Role:      msg.Role,
				Content:   msg.StringContent(),
				ToolCalls: msg.ToolCalls,
			})
		} else {
			mediaContents, err := msg.ParseContent()
			if err != nil {
				return nil, err
			}
			for _, mediaContent := range mediaContents {
				switch mediaContent.Type {
				case v1.ContentTypeText:
					ollamaReq.Messages = append(ollamaReq.Messages, OllamaMessage{
						Role:      msg.Role,
						Content:   mediaContent.Text,
						ToolCalls: msg.ToolCalls,
					})
					continue
				case v1.ContentTypeImageURL:
					if mediaContent.ImageUrl.Url == "" {
						continue
					}
					b64, err := tools.NewImageFileData(mediaContent.ImageUrl.Url, true)
					if err != nil {
						return nil, err
					}
					msg := OllamaMessage{
						Role:      msg.Role,
						Content:   mediaContent.ImageUrl.Url,
						Images:    []string{b64.URL},
						ToolCalls: msg.ToolCalls,
					}
					ollamaReq.Messages = append(ollamaReq.Messages, msg)
					continue
				default:
					continue
				}

			}
		}
	}
	reqBytes, err := sonic.Marshal(ollamaReq)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewBuffer(reqBytes)), nil
}

func (c *Client) convertChatCompletionsResponse(ollamaResp *OllamaChatCompletionResponse) *v1.ChatCompletionResponse {
	oaiResp := v1.ChatCompletionResponse{
		ID:                "chatcmpl-" + ulid.Make().String(),
		Model:             ollamaResp.Model,
		Object:            "chat.completion",
		Created:           ollamaResp.CreatedAt.Unix(),
		SystemFingerprint: "fp_" + ulid.Make().String(),
		Choices:           nil,
		Usage: &v1.Usage{
			PromptTokens:     ollamaResp.PromptEvalCount,
			CompletionTokens: ollamaResp.EvalCount,
			TotalTokens:      ollamaResp.EvalCount + ollamaResp.PromptEvalCount,
		},
	}
	choices := make([]v1.Choice, 0, len(ollamaResp.Message))
	for idx, msg := range ollamaResp.Message {
		oaiChoice := v1.Choice{
			Index: idx,
			Message: v1.CompletionMessage{
				Role:    msg.Role,
				Content: msg.Content,
			},
		}
		if idx == len(ollamaResp.Message)-1 {
			oaiChoice.FinishReason = "stop"
		}
		choices = append(choices, oaiChoice)
	}
	return &oaiResp
}

func (c *Client) convertChatCompletionsStreamResponse(ollamaResp *OllamaChatCompletionResponse) *v1.ChatCompletionStreamResponse {
	oaiResp := v1.ChatCompletionStreamResponse{
		ID:                "chatcmpl-" + ulid.Make().String(),
		Model:             ollamaResp.Model,
		Object:            "chat.completion.chunk",
		Created:           ollamaResp.CreatedAt.Unix(),
		SystemFingerprint: "fp_" + ulid.Make().String(),
		Choices:           nil,
		Usage:             nil,
	}
	choices := make([]v1.ChoiceWithDelta, 0, len(ollamaResp.Message))
	for idx, msg := range ollamaResp.Message {
		oaiChoice := v1.ChoiceWithDelta{
			Index: idx,
			Delta: v1.Delta{
				Role:      msg.Role,
				Content:   msg.Content,
				ToolCalls: msg.ToolCalls,
			},
		}
		if idx == len(ollamaResp.Message)-1 {
			oaiChoice.FinishReason = "stop"
		}
		choices = append(choices, oaiChoice)
	}
	if ollamaResp.Done {
		oaiResp.Usage = &v1.Usage{
			PromptTokens:     ollamaResp.PromptEvalCount,
			CompletionTokens: ollamaResp.EvalCount,
			TotalTokens:      ollamaResp.EvalCount + ollamaResp.PromptEvalCount,
		}
	}
	return &oaiResp
}

func (c *Client) CreateChatCompletions(ctx context.Context, req *v1.ChatCompletionRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := "/api/chat"
	reqBody, err := c.ConvertChatCompletions(req)
	if err != nil {
		return nil, nil, err
	}
	if req.Stream {
		body, header, err := base.Relay(ctx, http.MethodPost, targetUrl, reqBody, c.GenerateHeaderByContentType("application/json"), c.Client.Client)
		if err != nil {
			return nil, nil, err
		}
		var ollamaResp OllamaChatCompletionResponse
		err = sonic.ConfigDefault.NewDecoder(body).Decode(&ollamaResp)
		if err != nil {
			return nil, nil, err
		}
		oaiResp := c.convertChatCompletionsResponse(&ollamaResp)
		oaiRespBytes, err := sonic.Marshal(oaiResp)
		if err != nil {
			return nil, nil, err
		}
		return io.NopCloser(bytes.NewBuffer(oaiRespBytes)), header, nil
	} else {
		body, header, err := base.Relay(ctx, http.MethodPost, targetUrl, reqBody, c.GenerateHeaderByContentType("application/json"), c.Client.Client)
		if err != nil {
			return nil, nil, err
		}
		r, w := io.Pipe()
		go func(w *io.PipeWriter) {
			defer func() {
				err := body.Close()
				if err != nil {
					fmt.Println("Error closing body:", err)
				}
				err = w.Close()
				if err != nil {
					fmt.Println("Error closing w:", err)
				}
			}()
			scanner := bufio.NewScanner(body)
			for scanner.Scan() {
				line := scanner.Bytes()
				check := string(line)
				if chunk := strings.TrimSpace(check); len(chunk) == 0 {
					continue
				}
				var ollamaResp OllamaChatCompletionResponse
				err = sonic.Unmarshal(line, &ollamaResp)
				if err != nil {
					fmt.Println("Error decoding JSON:", err)
					continue
				}
				oaiResp := c.convertChatCompletionsStreamResponse(&ollamaResp)
				oaiRespBytes, err := sonic.Marshal(oaiResp)
				if err != nil {
					fmt.Printf("data:%#v", line)
					fmt.Println("sonic.Marshal失败", err)
					return
				}
				_, err = fmt.Fprintf(w, "data: %s\n\n", oaiRespBytes)
				if err != nil {
					fmt.Println("Error writing to pipe:", err)
					return
				}
			}
			_, _ = fmt.Fprintf(w, "data: [DONE]\n\n")
		}(w)
		return r, header, nil
	}
}

func (c *Client) ConvertCompletions(req *v1.CompletionsRequest) (io.ReadCloser, error) {
	var err error
	ollamaReq := OllamaCompletionRequest{
		Model:  req.Model,
		Prompt: req.Prompt,
		Stream: req.Stream,
		Suffix: req.Suffix,
		Images: req.Images,
		Format: req.Format,
		Options: &Options{
			NumCtx:           req.MaxTokens,
			Temperature:      req.Temperature,
			Seed:             req.Seed,
			Stop:             req.Stop,
			TopP:             req.TopP,
			PresencePenalty:  req.PresencePenalty,
			FrequencyPenalty: req.FrequencyPenalty,
		},
		System:    req.System,
		Template:  req.Template,
		Raw:       req.Raw,
		KeepAlive: req.KeepAlive,
	}
	bodyBytes, err := sonic.Marshal(ollamaReq)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewBuffer(bodyBytes)), nil
}

func (c *Client) convertCompletionsResponse(ollamaResp *OllamaCompletionResponse) *v1.CompletionsResp {
	resp := v1.CompletionsResp{
		ID:      "cmpl-" + ulid.Make().String(),
		Object:  "text_completion",
		Created: ollamaResp.CreatedAt.Unix(),
		Model:   ollamaResp.Model,
		Choices: nil,
		Usage:   nil,
	}
	choice := v1.CompletionsChoice{
		Text:  ollamaResp.Response,
		Index: 0,
	}
	if ollamaResp.Done {
		resp.Usage = &v1.Usage{
			CompletionTokens: ollamaResp.EvalCount,
			PromptTokens:     ollamaResp.PromptEvalCount,
			TotalTokens:      ollamaResp.EvalCount + ollamaResp.PromptEvalCount,
		}
		choice.FinishReason = "stop"
	}
	resp.Choices = []v1.CompletionsChoice{choice}
	return &resp
}

func (c *Client) CreateCompletions(ctx context.Context, req *v1.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := "/api/generate"
	reqBody, err := c.ConvertCompletions(req)
	if err != nil {
		return nil, nil, err
	}
	if req.Stream {
		body, header, err := base.Relay(ctx, http.MethodPost, targetUrl, reqBody, c.GenerateHeaderByContentType("application/json"), c.Client.Client)
		if err != nil {
			return nil, nil, err
		}
		var ollamaResp OllamaCompletionStreamResponse
		err = sonic.ConfigDefault.NewDecoder(body).Decode(&ollamaResp)
		if err != nil {
			return nil, nil, err
		}
		oaiResp := c.convertCompletionsResponse(&ollamaResp)
		oaiRespBytes, err := sonic.Marshal(oaiResp)
		if err != nil {
			return nil, nil, err
		}
		return io.NopCloser(bytes.NewBuffer(oaiRespBytes)), header, nil
	} else {
		body, header, err := base.Relay(ctx, http.MethodPost, targetUrl, reqBody, c.GenerateHeaderByContentType("application/json"), c.Client.Client)
		if err != nil {
			return nil, nil, err
		}
		r, w := io.Pipe()
		go func(w *io.PipeWriter) {
			defer func() {
				err := body.Close()
				if err != nil {
					fmt.Println("Error closing body:", err)
				}
				err = w.Close()
				if err != nil {
					fmt.Println("Error closing w:", err)
				}
			}()
			scanner := bufio.NewScanner(body)
			for scanner.Scan() {
				line := scanner.Bytes()
				check := string(line)
				if chunk := strings.TrimSpace(check); len(chunk) == 0 {
					continue
				}
				var ollamaResp OllamaCompletionStreamResponse
				err = sonic.Unmarshal(line, &ollamaResp)
				if err != nil {
					fmt.Println("Error decoding JSON:", err)
					continue
				}
				oaiResp := c.convertCompletionsResponse(&ollamaResp)
				oaiRespBytes, err := sonic.Marshal(oaiResp)
				if err != nil {
					fmt.Printf("data:%#v", line)
					fmt.Println("sonic.Marshal失败", err)
					return
				}
				_, err = fmt.Fprintf(w, "data: %s\n\n", oaiRespBytes)
				if err != nil {
					fmt.Println("Error writing to pipe:", err)
					return
				}
			}
			_, _ = fmt.Fprintf(w, "data: [DONE]\n\n")
		}(w)
		return r, header, nil
	}
}

func (c *Client) Models(ctx context.Context) (*v1.ModelResponse, error) {
	var err error
	targetUrl := c.EndPoint + "/v1/models"
	data, _, err := base.Relay(ctx, http.MethodGet, targetUrl, nil, nil, c.Client.Client)
	if err != nil {
		return nil, err
	}
	dataBytes, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("read all error: %w", err)
	}
	var resp v1.ModelResponse
	err = sonic.Unmarshal(dataBytes, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return &resp, nil
}

func (c *Client) CreateEmbeddings(ctx context.Context, req *v1.EmbeddingsRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := "/api/embed"
	ollamaReq := OllamaEmbeddingRequest{
		Model:     req.Model,
		Input:     req.Input,
		Options:   &Options{},
		KeepAlive: "",
	}
	reqBytes, err := sonic.Marshal(ollamaReq)
	if err != nil {
		return nil, nil, err
	}
	resp, header, err := base.Relay(ctx, http.MethodPost, targetUrl, bytes.NewBuffer(reqBytes), c.GenerateHeaderByContentType("application/json"), c.Client.Client)
	if err != nil {
		return nil, nil, err
	}
	respBytes, err := io.ReadAll(resp)
	if err != nil {
		return nil, nil, err
	}
	var ollamaResp OllamaEmbeddingResponse
	err = sonic.Unmarshal(respBytes, &ollamaResp)
	if err != nil {
		return nil, nil, err
	}
	oaiResp := v1.EmbeddingsResponse{
		Object: "list",
		Data:   make([]v1.EmbeddingsData, 0, len(ollamaResp.Embeddings)),
		Usage: v1.Usage{
			PromptTokens:     ollamaResp.PromptEvalCount,
			CompletionTokens: ollamaResp.PromptEvalCount,
			TotalTokens:      ollamaResp.PromptEvalCount,
		},
	}
	for _, embedding := range ollamaResp.Embeddings {
		oaiResp.Data = append(oaiResp.Data, v1.EmbeddingsData{
			ID:        "",
			Object:    "embedding",
			Index:     0,
			Embedding: embedding,
		})
	}
	oaiRespBytes, err := sonic.Marshal(oaiResp)
	if err != nil {
		return nil, nil, err
	}
	return io.NopCloser(bytes.NewBuffer(oaiRespBytes)), header, nil
}

func (c *Client) CreateRerank(ctx context.Context, req *v1.RerankRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateSpeech(ctx context.Context, req *v1.AudioSpeechRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateTranslation(ctx context.Context, req *v1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateTranscription(ctx context.Context, req *v1.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateImage(ctx context.Context, req *v1.ImageGenerateRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateImageEdit(ctx context.Context, req *v1.ImageEditRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateImageVariation(ctx context.Context, req *v1.ImageVariationRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateVideoSubmit(ctx context.Context, req *v1.VideoRequest) (*v1.VideoResponse, error) {
	_, _, err := base.NoImplementMethod(ctx, req)
	return nil, err
}

func (c *Client) GetVideoStatus(ctx context.Context, externalID string) (bool, any, error) {
	_, _, err := base.NoImplementMethod(ctx, externalID)
	return false, nil, err
}
