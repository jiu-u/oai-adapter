package api

import (
	"encoding/json"
	"fmt"
)

type (
	ChatRequest struct {
		Messages            []Message       `json:"messages"`                 // 聊天消息
		Model               string          `json:"model" binding:"required"` // 模型名称
		Store               bool            `json:"store,omitempty"`
		ReasoningEffect     string          `json:"reasoning_effect,omitempty"`
		MetaData            any             `json:"meta_data,omitempty"`
		FrequencyPenalty    float64         `json:"frequency_penalty,omitempty"` // 频率惩罚
		LogitBias           any             `json:"logit_bias,omitempty"`
		Logprobs            bool            `json:"logprobs,omitempty"`
		TopLogprobs         int             `json:"top_logprobs,omitempty"`
		MaxTokens           int             `json:"max_tokens,omitempty"`     // 最大 token 数
		MaxCompletionTokens int             `json:"max_completion,omitempty"` // 最大完成数
		N                   int             `json:"n,omitempty"`
		Modelities          any             `json:"modelities,omitempty"`
		Prediction          any             `json:"prediction,omitempty"`
		Audio               any             `json:"audio,omitempty"`
		PresencePenalty     float64         `json:"presence_penalty,omitempty"` // 存在惩罚
		ResponseFormat      *ResponseFormat `json:"response_format,omitempty"`
		Seed                int             `json:"seed,omitempty"`
		ServiceTier         string          `json:"service_tier,omitempty"`
		Stop                any             `json:"stop,omitempty"`          // 停止标志
		Stream              bool            `json:"stream"`                  // 是否流式返回
		SteamOptions        *StreamOptions  `json:"steam_options,omitempty"` // 流式配置
		Temperature         float64         `json:"temperature,omitempty"`   // 温度
		TopP                float64         `json:"top_p,omitempty"`         // top_p
		Tools               []ToolCall      `json:"tools,omitempty"`
		ToolChoice          any             `json:"tool_choice,omitempty"`
		ParallelToolCalls   bool            `json:"parallel_tool_calls,omitempty"`
		User                string          `json:"user,omitempty"`
		FunctionCall        any             `json:"function_call,omitempty"`
		Functions           any             `json:"functions,omitempty"`
	}
	ResponseFormat struct {
		Type       string            `json:"type,omitempty"`
		JsonSchema *FormatJsonSchema `json:"json_schema,omitempty"`
	}
	FormatJsonSchema struct {
		Description string `json:"description,omitempty"`
		Name        string `json:"name"`
		Schema      any    `json:"schema,omitempty"`
		Strict      any    `json:"strict,omitempty"`
	}
	StreamOptions struct {
		IncludeUsage bool `json:"include_usage,omitempty"`
	}
)

type Message struct {
	Role         string          `json:"role"`
	Content      json.RawMessage `json:"content"`
	Name         string          `json:"name,omitempty"`
	ToolCallId   string          `json:"tool_call_id,omitempty"`
	Refusal      any             `json:"refusal,omitempty"`
	Audio        any             `json:"audio,omitempty"`
	ToolCalls    []ToolCall      `json:"tool_calls,omitempty"`
	FunctionCall any             `json:"function_call,omitempty"`
}

type MediaContent struct {
	Type       string             `json:"type"`
	Text       string             `json:"text"`
	Refusal    string             `json:"refusal,omitempty"`
	ImageUrl   *MessageImageUrl   `json:"image_url,omitempty"`
	InputAudio *MessageInputAudio `json:"input_audio,omitempty"`
}

type MessageImageUrl struct {
	Url    string `json:"url"`
	Detail string `json:"detail"`
}

type MessageInputAudio struct {
	Data   string `json:"data"` //base64
	Format string `json:"format"`
}

const (
	ContentTypeText       = "text"
	ContentTypeImageURL   = "image_url"
	ContentTypeInputAudio = "input_audio"
)

func (m *Message) StringContent() string {
	var stringContent string
	if err := json.Unmarshal(m.Content, &stringContent); err == nil {
		return stringContent
	}
	return string(m.Content)
}

func (m *Message) SetStringContent(content string) {
	jsonContent, _ := json.Marshal(content)
	m.Content = jsonContent
}

func (m *Message) IsStringContent() bool {
	var stringContent string
	if err := json.Unmarshal(m.Content, &stringContent); err == nil {
		return true
	}
	return false
}

func (m *Message) ParseContent() ([]MediaContent, error) {
	var contentList []MediaContent
	var stringContent string
	if err := json.Unmarshal(m.Content, &stringContent); err == nil {
		contentList = append(contentList, MediaContent{
			Type: ContentTypeText,
			Text: stringContent,
		})
		return contentList, nil
	}

	if err := json.Unmarshal(m.Content, &contentList); err == nil {
		return contentList, nil
	} else {
		return nil, fmt.Errorf("parse content error: %w", err)

	}
}
