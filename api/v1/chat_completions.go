package v1

import (
	"encoding/json"
	"fmt"
)

type (
	ChatCompletionRequest struct {
		Messages         []Message `json:"messages"`                 // 聊天消息
		Model            string    `json:"model" binding:"required"` // 模型名称
		Audio            any       `json:"audio,omitempty"`
		FrequencyPenalty float64   `json:"frequency_penalty,omitempty" binding:"omitempty,min=-2,max=2" default:"0"` // 频率惩罚
		// Deprecated: Use tool_choice instead.
		FunctionCall any `json:"function_call,omitempty"` //string(auto,none) or object{"name":"fn_name"}
		// Deprecated: Use tools instead.
		Functions           Function `json:"functions,omitempty"`
		LogitBias           any      `json:"logit_bias,omitempty"`
		Logprobs            bool     `json:"logprobs,omitempty"`
		MaxCompletionTokens int      `json:"max_completion_tokens,omitempty"` // 最大完成数
		// Deprecated: Use max_completion_tokens instead.
		MaxTokens         int             `json:"max_tokens,omitempty"`
		MetaData          any             `json:"meta_data,omitempty"`
		Modalities        []string        `json:"modalities,omitempty"`
		N                 int             `json:"n,omitempty" binding:"omitempty" default:"1"`
		ParallelToolCalls bool            `json:"parallel_tool_calls,omitempty"`
		Prediction        any             `json:"prediction,omitempty"`
		PresencePenalty   float64         `json:"presence_penalty,omitempty" binding:"omitempty,min=-2,max=2" default:"0"` // 存在惩罚
		ReasoningEffect   string          `json:"reasoning_effect,omitempty" binding:"omitempty,oneof=low medium high"`
		ResponseFormat    *ResponseFormat `json:"response_format,omitempty"`
		Seed              int64           `json:"seed,omitempty"`
		ServiceTier       string          `json:"service_tier,omitempty"`
		Stop              any             `json:"stop,omitempty"` // 停止标志
		Store             bool            `json:"store,omitempty"`
		Stream            bool            `json:"stream" binding:"omitempty" default:"false"` // 流式返回
		StreamOptions     *StreamOptions  `json:"stream_options,omitempty"`
		Temperature       float64         `json:"temperature,omitempty" binding:"omitempty,min=0,max=2" default:"1"` // 温度
		ToolChoice        any             `json:"tool_choice,omitempty"`
		Tools             []Tool          `json:"tools,omitempty"`
		TopLogprobs       int             `json:"top_logprobs,omitempty" binding:"omitempty,min=0,max=20"`
		TopP              float64         `json:"top_p,omitempty"`
		User              string          `json:"user,omitempty"`
		WebSearchOptions  any             `json:"web_search_options,omitempty"`
	}
)

type Message struct {
	Content json.RawMessage `json:"content"`
	Role    string          `json:"role"`
	Name    string          `json:"name,omitempty"`
	//
	Audio   any `json:"audio,omitempty"`
	Refusal any `json:"refusal,omitempty"` // string or null
	// Deprecated: Use tool_calls instead.
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
	ToolCalls    []ToolCall    `json:"tool_calls,omitempty"`
	ToolCallId   string        `json:"tool_call_id,omitempty"`
}

type MediaContent struct {
	Type       string      `json:"type"`
	Text       string      `json:"text"`
	Refusal    string      `json:"refusal,omitempty"`
	ImageUrl   *ImageUrl   `json:"image_url,omitempty"`
	InputAudio *InputAudio `json:"input_audio,omitempty"`
	File       *File       `json:"file,omitempty"`
}

type ImageUrl struct {
	Url    string `json:"url" binding:"required"`
	Detail string `json:"detail"`
}

type InputAudio struct {
	Data   string `json:"data"` //base64
	Format string `json:"format"`
}

type File struct {
	FileData any    `json:"file_data"`
	FileId   any    `json:"file_id"`
	Filename string `json:"filename"`
}

const (
	ContentTypeText       = "text"
	ContentTypeImageURL   = "image_url"
	ContentTypeInputAudio = "input_audio"
	ContentTypeFile       = "file"
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

func messageUseExample(msg *Message) {
	if msg.IsStringContent() {
		fmt.Println(msg.StringContent())
		return
	}
	mediaContents, err := msg.ParseContent()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, mediaContent := range mediaContents {
		switch mediaContent.Type {
		case ContentTypeText:
			fmt.Println(mediaContent.Text)
			continue
		case ContentTypeImageURL:
			if mediaContent.ImageUrl.Url == "" {
				continue
			}
			fmt.Println(mediaContent.ImageUrl.Url)
			continue
		case ContentTypeInputAudio:
			if mediaContent.InputAudio.Data == "" {
				continue
			}
			if mediaContent.InputAudio.Format == "" {
				continue
			}
			fmt.Println(mediaContent.InputAudio.Data)
			fmt.Println(mediaContent.InputAudio.Format)
			continue
		case ContentTypeFile:
			if mediaContent.File.Filename == "" {
				continue
			}
			fmt.Println(mediaContent.File.Filename)
			continue
		default:
			fmt.Println("unknown content type")
			continue
		}
	}

}
