package v1

import "encoding/json"

type (
	Function struct {
		Name        string          `json:"name" binding:"required"`
		Description string          `json:"description,omitempty"`
		Parameters  json.RawMessage `json:"parameters,omitempty"`
		// tool_calls
		Arguments string `json:"arguments,omitempty"`
	}
	FunctionCall struct {
		Name      string `json:"name" binding:"required"`
		Arguments string `json:"arguments,omitempty"`
	}
)

type (
	ToolChoice struct {
		Type     string          `json:"type" binding:"required"`
		Function *FunctionChoice `json:"function,omitempty"`
	}
	FunctionChoice struct {
		Name string `json:"name" binding:"required"`
	}
	Tool struct {
		Type     string   `json:"type" binding:"required"`
		Function Function `json:"function"`
	}
	ToolCall struct {
		Id       string   `json:"id,omitempty"`
		Type     string   `json:"type"`
		Function Function `json:"function"`
		// streaming
		Index int `json:"index,omitempty"`
	}
	WebSearchOptions struct {
		SearchContentSize string        `json:"search_content_size,omitempty" binding:"omitempty,oneof=low medium high" default:"medium"`
		UserLocation      *UserLocation `json:"user_location,omitempty"`
	}
	UserLocation struct {
		Approximate *Approximate `json:"approximate,omitempty"`
		Type        string       `json:"type,omitempty"`
	}
	Approximate struct {
		City     string `json:"city,omitempty"`
		Country  string `json:"country,omitempty"`
		Region   string `json:"region,omitempty"`
		Timezone string `json:"timezone,omitempty"`
	}
)

type (
	ResponseFormat struct {
		Type       string            `json:"type,omitempty"`
		JsonSchema *FormatJsonSchema `json:"json_schema,omitempty"`
	}
	FormatJsonSchema struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description,omitempty"`
		Schema      any    `json:"schema,omitempty"`
		Strict      bool   `json:"strict,omitempty" binding:"omitempty" default:"false"`
	}
	StreamOptions struct {
		IncludeUsage bool `json:"include_usage,omitempty"`
	}
	Prediction struct {
		Type    string          `json:"type"`
		Content json.RawMessage `json:"content"` // string or array(ContentPart)
	}
	ContentPart struct {
		Text string `json:"text" binding:"required"`
		Type string `json:"type" binding:"required"`
	}
	Audio struct {
		Id string `json:"id"`
		// completion response message
		Data       string `json:"data,omitempty"`
		ExpiresAt  int64  `json:"expires_at,omitempty"`
		Transcript string `json:"transcript,omitempty"`
	}
	Reasoning struct {
		Effect          string `json:"effect,omitempty"`
		GenerateSummary string `json:"generate_summary,omitempty"`
		Summary         string `json:"summary,omitempty"`
	}
)

type (
	Usage struct {
		CompletionTokens        int                      `json:"completion_tokens"`
		PromptTokens            int                      `json:"prompt_tokens"`
		TotalTokens             int                      `json:"total_tokens"`
		CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
		PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`

		InputTokens        int                 `json:"input_tokens,omitempty"`
		OutputTokens       int                 `json:"output_tokens,omitempty"`
		InputTokenDetails  *InputTokenDetails  `json:"input_token_details,omitempty"`
		OutputTokenDetails *OutputTokenDetails `json:"output_token_details,omitempty"`
	}
	CompletionTokensDetails struct {
		AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"`
		AudioTokens              int `json:"audio_tokens,omitempty"`
		ReasoningTokens          int `json:"reasoning_tokens,omitempty"`
		RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"`
	}
	PromptTokensDetails struct {
		AudioTokens  int `json:"audio_tokens,omitempty"`
		CachedTokens int `json:"cached_tokens,omitempty"`
	}
	InputTokenDetails struct {
		CachedTokens         int `json:"cached_tokens"`
		CachedCreationTokens int `json:"-"`
		TextTokens           int `json:"text_tokens"`
		AudioTokens          int `json:"audio_tokens"`
		ImageTokens          int `json:"image_tokens"`
	}
	OutputTokenDetails struct {
		TextTokens      int `json:"text_tokens"`
		AudioTokens     int `json:"audio_tokens"`
		ReasoningTokens int `json:"reasoning_tokens"`
	}
)
