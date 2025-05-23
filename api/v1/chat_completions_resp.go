package v1

type (
	ChatCompletionResponse struct {
		ID                string   `json:"id"`
		Model             string   `json:"model"`
		Object            string   `json:"object"`
		Created           int64    `json:"created"`
		ServiceTier       string   `json:"service_tier,omitempty"`
		SystemFingerprint string   `json:"systemFingerprint,omitempty"`
		Choices           []Choice `json:"choices"`
		Usage             Usage    `json:"usage,omitempty"`
	}
	Choice struct {
		Index        int               `json:"index"`
		Message      CompletionMessage `json:"message"`
		FinishReason string            `json:"finish_reason"`
		Logprobs     any               `json:"logprobs,omitempty"`
	}
	CompletionMessage struct {
		Content      string        `json:"content,omitempty"`
		Refusal      any           `json:"refusal,omitempty"` // string or null
		Role         string        `json:"role"`
		Annotations  []Annotation  `json:"annotations,omitempty"`
		Audio        *Audio        `json:"audio,omitempty"`
		FunctionCall *FunctionCall `json:"function_call,omitempty"`
		ToolCalls    []ToolCall    `json:"tool_calls,omitempty"`
	}
	Annotation struct {
		Type        string `json:"type"`
		UrlCitation any    `json:"url_citation,omitempty"`
	}
)

// streaming
type (
	ChatCompletionStreamResponse struct {
		Choices           []ChoiceWithDelta `json:"choices"`
		Created           int64             `json:"created"`
		ID                string            `json:"id"`
		Model             string            `json:"model"`
		Object            string            `json:"object"`
		SystemFingerprint string            `json:"systemFingerprint,omitempty"`
		Usage             *Usage            `json:"usage,omitempty"`
	}

	ChoiceWithDelta struct {
		Delta        Delta  `json:"delta"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
		Logprobs     any    `json:"logprobs,omitempty"`
	}
	Delta struct {
		Role    string `json:"role"`
		Content string `json:"content"`
		// Deprecated
		FunctionCall *FunctionCall `json:"function_call,omitempty"`
		Refusal      string        `json:"refusal,omitempty"`
		ToolCalls    []ToolCall    `json:"tool_calls,omitempty"`
	}
)
