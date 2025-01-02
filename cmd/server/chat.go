package main

import (
	"encoding/json"
	"github.com/jiu-u/oai-adapter/api"
	"github.com/jiu-u/oai-adapter/clients/gemini"
	"net/http"
)

func HandleChatCompletions(cl *gemini.Client) func(w http.ResponseWriter, r *http.Request) {
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
