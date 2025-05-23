package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	oaiadapter "github.com/jiu-u/oai-adapter"
	v1 "github.com/jiu-u/oai-adapter/api/v1"
	"github.com/jiu-u/oai-adapter/common"
	"io"
	"net/http"
	"time"
)

type RelayAction int

const (
	RelayRequest RelayAction = iota
	Responses    RelayAction = iota
	ChatCompletions
	Completions
	Embeddings
	Rerank
	CreateSpeech
	Transcriptions
	Translations
	CreateImage
	CreateImageEdit
	ImageVariations
	VideoSubmit
	VideoStatus
)

func RelayHandler(cl oaiadapter.Adapter, action RelayAction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		//	读取数据
		if r.Method == http.MethodPost {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			//fmt.Println(string(data))
			r.Body = io.NopCloser(bytes.NewBuffer(data))
		}
		fmt.Printf("函数执行耗时1: %s\n", time.Since(startTime))
		requestBody, err := ParseRequest(r, action)
		fmt.Printf("函数执行耗时2: %s\n", time.Since(startTime))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		respBody, respHeader, err := DoRelayRequest(r.Context(), cl, action, requestBody)

		elapsed := time.Since(startTime)
		fmt.Printf("函数执行耗时: %s\n", elapsed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		HandleOAIResponse(w, r, respBody, respHeader)
	}
}

func ParseRequest(r *http.Request, action RelayAction) (any, error) {
	switch action {
	case Responses:
		var req v1.ResponsesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case ChatCompletions:
		var req v1.ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case Completions:
		var req v1.CompletionsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case Embeddings:
		var req v1.EmbeddingsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case Rerank:
		var req v1.RerankRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case CreateSpeech:
		var req v1.AudioSpeechRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case Transcriptions:
		var req v1.TranscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case Translations:
		var req v1.TranslationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case CreateImage:
		var req v1.ImageGenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case CreateImageEdit:
		var req v1.ImageEditRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case ImageVariations:
		var req v1.ImageVariationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	case VideoSubmit:
		var req v1.VideoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	default:
		return nil, fmt.Errorf("unsupported action")
	}
}

func DoRelayRequest(ctx context.Context, cl oaiadapter.Adapter, action RelayAction, requestBody any) (io.ReadCloser, http.Header, error) {
	switch action {
	case Responses:
		req := requestBody.(*v1.ResponsesRequest)
		return cl.CreateResponses(ctx, req)
	case ChatCompletions:
		req := requestBody.(*v1.ChatCompletionRequest)
		return cl.CreateChatCompletions(ctx, req)
	case Completions:
		req := requestBody.(*v1.CompletionsRequest)
		return cl.CreateCompletions(ctx, req)
	case Embeddings:
		req := requestBody.(*v1.EmbeddingsRequest)
		return cl.CreateEmbeddings(ctx, req)
	case Rerank:
		req := requestBody.(*v1.RerankRequest)
		return cl.CreateRerank(ctx, req)
	case CreateSpeech:
		req := requestBody.(*v1.AudioSpeechRequest)
		return cl.CreateSpeech(ctx, req)
	case Transcriptions:
		req := requestBody.(*v1.TranscriptionRequest)
		return cl.CreateTranscription(ctx, req)
	case Translations:
		req := requestBody.(*v1.TranslationRequest)
		return cl.CreateTranslation(ctx, req)
	case CreateImage:
		req := requestBody.(*v1.ImageGenerateRequest)
		return cl.CreateImage(ctx, req)
	case CreateImageEdit:
		req := requestBody.(*v1.ImageEditRequest)
		return cl.CreateImageEdit(ctx, req)
	case ImageVariations:
		req := requestBody.(*v1.ImageVariationRequest)
		return cl.CreateImageVariation(ctx, req)
	case VideoSubmit:
		req := requestBody.(*v1.VideoRequest)
		data, err := cl.CreateVideoSubmit(ctx, req)
		if err != nil {
			return nil, nil, err
		}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, nil, err
		}
		header := http.Header{}
		header.Set("Content-Type", "application/json")
		header.Set("Cache-Control", "no-cache")
		return io.NopCloser(bytes.NewBuffer(dataBytes)), header, nil
	default:
		return nil, nil, fmt.Errorf("unsupported action")
	}
}

func HandleModels(cl oaiadapter.Adapter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := cl.Models(r.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleVideoStatus(w http.ResponseWriter, r *http.Request) {
	var query v1.VideoStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if query.RequestId == "" {
		http.Error(w, "requestId is empty", http.StatusBadRequest)
		return
	}
	taskManager := common.GetDefaultTaskManager()
	if taskManager == nil {
		http.Error(w, "taskManager is nil", http.StatusInternalServerError)
		return
	}
	res, exist := taskManager.GetTaskResult(query.RequestId)
	if !exist {
		http.Error(w, "task not exist", http.StatusBadRequest)
		return
	}
	resData := res.Data.(*v1.VideoStatusResponse)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
