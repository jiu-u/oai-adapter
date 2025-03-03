package ollama

import "time"

type (
	OllamaGenerateRequest struct {
		Model     string   `json:"model"`
		Prompt    string   `json:"prompt"`
		Suffix    string   `json:"suffix,omitempty"`
		Images    []string `json:"images,omitempty"`
		Format    any      `json:"format,omitempty"`
		Options   Options  `json:"options,omitempty"`
		System    string   `json:"system,omitempty"`
		Template  string   `json:"template,omitempty"`
		Raw       bool     `json:"raw,omitempty"`
		KeepAlive string   `json:"keep_alive,omitempty"`
		Stream    bool     `json:"stream,omitempty"`
	}
	Options struct {
		Temperature      float64 `json:"temperature,omitempty"`
		Seed             any     `json:"seed,omitempty"`
		Echo             bool    `json:"echo,omitempty"`
		FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
		Logprobs         any     `json:"logprobs,omitempty"`
		MaxTokens        int64   `json:"max_tokens,omitempty"`
		BestOf           int64   `json:"best_of,omitempty"`
		LogitBias        any     `json:"logit_bias,omitempty"`
		N                any     `json:"n,omitempty"`
		PresencePenalty  float64 `json:"presence_penalty,omitempty"`
		Stop             any     `json:"stop,omitempty"`
		TopP             float64 `json:"top_p,omitempty"`
		User             string  `json:"user,omitempty"`
	}
	OllamaGenerateResponse struct {
		Model              string    `json:"model"`
		CreatedAt          time.Time `json:"created_at"`
		Response           string    `json:"response"`
		Done               bool      `json:"done"`
		DoneReason         string    `json:"done_reason"`
		Context            []int     `json:"context"`
		TotalDuration      int64     `json:"total_duration"`
		LoadDuration       int64     `json:"load_duration"`
		PromptEvalCount    int       `json:"prompt_eval_count"`
		PromptEvalDuration int64     `json:"prompt_eval_duration"`
		EvalCount          int       `json:"eval_count"`
		EvalDuration       int64     `json:"eval_duration"`
	}
)
