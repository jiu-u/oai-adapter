package main

import (
	"context"
	"encoding/json"
	"fmt"
	oaiadapter "github.com/jiu-u/oai-adapter"
	"github.com/jiu-u/oai-adapter/api"
	"io"
	"net/http"
)

type RelayAction int

const (
	ChatCompletions RelayAction = iota
	ChatCompletionsByBytes
	Completions
	CompletionsByBytes
	Embeddings
	EmbeddingsByBytes
	CreateSpeech
	CreateSpeechByBytes
	Transcriptions
	Translations
	CreateImage
	CreateImageByBytes
	CreateImageEdit
	ImageVariations
)

func RelayHandler(cl oaiadapter.Adapter, action RelayAction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ParseRequest(r, action)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		respBody, respHeader, err := DoRelayRequest(r.Context(), cl, action, requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		HandleOAIResponse(w, r, respBody, respHeader)
	}
}

func ParseRequest(r *http.Request, action RelayAction) (any, error) {
	switch action {
	case ChatCompletions:
		var req api.ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case Completions:
		var req api.CompletionsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case Embeddings:
		var req api.EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case CreateSpeech:
		var req api.SpeechRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case Transcriptions:
		var req api.TranscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case Translations:
		var req api.TranslationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case CreateImage:
		var req api.CreateImageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case CreateImageEdit:
		var req api.EditImageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case ImageVariations:
		var req api.CreateImageVariationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case ChatCompletionsByBytes, CompletionsByBytes, EmbeddingsByBytes, CreateSpeechByBytes, CreateImageByBytes:
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	default:
		return nil, fmt.Errorf("unsupported action")
	}
}

func DoRelayRequest(ctx context.Context, cl oaiadapter.Adapter, action RelayAction, requestBody any) (io.ReadCloser, http.Header, error) {
	switch action {
	case ChatCompletions:
		req := requestBody.(*api.ChatRequest)
		return cl.ChatCompletions(ctx, req)
	case ChatCompletionsByBytes:
		req := requestBody.([]byte)
		return cl.ChatCompletionsByBytes(ctx, req)
	case Completions:
		req := requestBody.(*api.CompletionsRequest)
		return cl.Completions(ctx, req)
	case CompletionsByBytes:
		req := requestBody.([]byte)
		return cl.CompletionsByBytes(ctx, req)
	case Embeddings:
		req := requestBody.(*api.EmbeddingRequest)
		return cl.Embeddings(ctx, req)
	case EmbeddingsByBytes:
		req := requestBody.([]byte)
		return cl.EmbeddingsByBytes(ctx, req)
	case CreateSpeech:
		req := requestBody.(*api.SpeechRequest)
		return cl.CreateSpeech(ctx, req)
	case CreateSpeechByBytes:
		req := requestBody.([]byte)
		return cl.CreateSpeechByBytes(ctx, req)
	case Transcriptions:
		req := requestBody.(*api.TranscriptionRequest)
		return cl.Transcriptions(ctx, req)
	case Translations:
		req := requestBody.(*api.TranslationRequest)
		return cl.Translations(ctx, req)
	case CreateImage:
		req := requestBody.(*api.CreateImageRequest)
		return cl.CreateImage(ctx, req)
	case CreateImageByBytes:
		req := requestBody.([]byte)
		return cl.CreateImageByBytes(ctx, req)
	case CreateImageEdit:
		req := requestBody.(*api.EditImageRequest)
		return cl.CreateImageEdit(ctx, req)
	case ImageVariations:
		req := requestBody.(*api.CreateImageVariationRequest)
		return cl.ImageVariations(ctx, req)
	default:
		return nil, nil, fmt.Errorf("unsupported action")
	}
}

func HandleChatCompletions(cl oaiadapter.Adapter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req api.ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//req.Model = "gemini-1.5-flash"
		resp, header, err := cl.ChatCompletions(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		HandleOAIResponse(w, r, resp, header)
	}
}

type ModelItem struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created,omitempty"`
	OwnedBy string `json:"owned_by,omitempty"`
}

type ModelResponse struct {
	Object string      `json:"object"`
	Data   []ModelItem `json:"data"`
}

func HandleModels(cl oaiadapter.Adapter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HandleModels")
		list, err := cl.Models(r.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		models := make([]ModelItem, len(list))
		fmt.Println(list)
		for i, model := range list {
			models[i] = ModelItem{
				ID:      model,
				Object:  "model",
				Created: 0,
				OwnedBy: "",
			}
		}
		resp := ModelResponse{
			Object: "list",
			Data:   models,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleCompletions(cl oaiadapter.Adapter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req api.CompletionsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, header, err := cl.Completions(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		HandleOAIResponse(w, r, resp, header)
	}
}

func HandleEmbeddings(cl oaiadapter.Adapter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req api.EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, header, err := cl.Embeddings(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		HandleOAIResponse(w, r, resp, header)
	}
}
