package api

type (
	ChatCompletionNoStreamResponse struct {
		ID                string   `json:"id"`
		Object            string   `json:"object"`
		Created           int64    `json:"created"`
		Model             string   `json:"model"`
		ServiceTier       string   `json:"service_tier,omitempty"`
		SystemFingerprint string   `json:"systemFingerprint,omitempty"`
		Choices           []Choice `json:"choices"`
		Usage             Usage    `json:"usage,omitempty"`
	}
	ChatCompletionStreamResponse struct {
		ID      string            `json:"id"`
		Object  string            `json:"object"`
		Created int64             `json:"created"`
		Model   string            `json:"model"`
		Choices []ChoiceWithDelta `json:"choices"`
		Usage   Usage             `json:"usage,omitempty"`
	}
	Choice struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	}

	ChoiceWithDelta struct {
		Index        int    `json:"index"`
		Delta        Delta  `json:"delta"`
		FinishReason string `json:"finish_reason"`
	}
	Delta struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	Usage struct {
		CompletionTokens        int                      `json:"completion_tokens"`
		PromptTokens            int                      `json:"prompt_tokens"`
		TotalTokens             int                      `json:"total_tokens"`
		CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
		PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`
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

	ChoiceMessage struct{}
)
