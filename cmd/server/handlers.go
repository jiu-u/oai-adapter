package main

import (
	"encoding/json"
	"fmt"
	oaiadapter "github.com/jiu-u/oai-adapter"
	"github.com/jiu-u/oai-adapter/api"
	"net/http"
)

func HandleChatCompletions(cl oaiadapter.Adapter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req api.ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.Model = "gemini-1.5-flash"
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
