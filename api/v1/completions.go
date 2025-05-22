package v1

// Legacy API

type CompletionsRequest struct {
	Model            string         `json:"model" binding:"required"`
	Prompt           string         `json:"prompt" binding:"required"`
	BestOf           int64          `json:"best_of,omitempty"`
	Echo             bool           `json:"echo,omitempty"`
	FrequencyPenalty float64        `json:"frequency_penalty,omitempty"`
	LogitBias        any            `json:"logit_bias,omitempty"`
	Logprobs         any            `json:"logprobs,omitempty"`
	MaxTokens        int64          `json:"max_tokens,omitempty"`
	N                any            `json:"n,omitempty"`
	PresencePenalty  float64        `json:"presence_penalty,omitempty"`
	Seed             any            `json:"seed,omitempty"`
	Stop             any            `json:"stop,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	StreamOptions    *StreamOptions `json:"stream_options,omitempty"`
	Suffix           string         `json:"suffix,omitempty"`
	Temperature      float64        `json:"temperature,omitempty"`
	TopP             float64        `json:"top_p,omitempty"`
	User             string         `json:"user,omitempty"`
	// ollama specific
	Images    []string `json:"images,omitempty"`
	Format    any      `json:"format,omitempty"`
	System    string   `json:"system,omitempty"`
	Template  string   `json:"template,omitempty"`
	Raw       bool     `json:"raw,omitempty"`
	KeepAlive string   `json:"keep_alive,omitempty"`
}

type (
	CompletionsResp struct {
		ID                string              `json:"id"`
		Object            string              `json:"object"`
		Created           int64               `json:"created"`
		Model             string              `json:"model"`
		Choices           []CompletionsChoice `json:"choices"`
		Usage             Usage               `json:"usage,omitempty"`
		SystemFingerprint string              `json:"systemFingerprint,omitempty"`
	}
	CompletionsChoice struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		Logprobs     any    `json:"logprobs,omitempty"`
		FinishReason string `json:"finish_reason,omitempty"`
	}
)
