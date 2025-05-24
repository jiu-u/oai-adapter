package fal

import "encoding/json"

type QueueResponse struct {
	Status      string          `json:"status"`
	RequestId   string          `json:"request_id"`
	ResponseUrl string          `json:"response_url"`
	StatusUrl   string          `json:"status_url"`
	CancelUrl   string          `json:"cancel_url"`
	Logs        json.RawMessage `json:"logs"`
	Metrics     json.RawMessage `json:"metrics"`
}

type (
	LLMRequest struct {
		Model        string `json:"model"`
		Prompt       string `json:"prompt" binding:"required"`
		SystemPrompt string `json:"system_prompt,omitempty"`
		Reasoning    bool   `json:"reasoning"`
	}
	LLMResponse struct {
		Output    string `json:"output"`
		Reasoning string `json:"reasoning"` // 使用指针类型表示可能为null
		Partial   bool   `json:"partial"`
		Error     string `json:"error"` // 使用指针类型表示可能为null
	}
)

type (
	ImageCreateRequest struct {
		Prompt              string     `json:"prompt" binding:"required"`
		ImageSize           *ImageSize `json:"image_size,omitempty"`
		NegativePrompt      string     `json:"negative_prompt,omitempty"`
		Seed                int        `json:"seed,omitempty"`
		NumInferenceSteps   int        `json:"num_inference_steps,omitempty"`
		NumImages           int        `json:"num_images,omitempty"`
		EnableSafetyChecker bool       `json:"enable_safety_checker"`
		SafetyTolerance     float64    `json:"safety_tolerance"`
	}
	ImageSize struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
	ImageCreateResponse struct {
		Images []Image `json:"images"`
		Prompt string  `json:"prompt"`
	}
	Image struct {
		URL         string `json:"url"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		ContentType string `json:"content_type"`
	}
)

type (
	Text2AudioRequest struct {
		Prompt string  `json:"prompt"`
		Voice  string  `json:"voice"`
		Speed  float64 `json:"speed,omitempty"`
	}
	Text2AudioResponse struct {
		Audio Audio `json:"audio"`
	}
	Audio struct {
		URL string `json:"url"`
	}
)

type (
	TextToVideoRequest struct {
		Prompt         string `json:"prompt,omitempty"`
		AspectRatio    string `json:"aspect_ratio,omitempty"`
		Resolution     string `json:"resolution,omitempty"`
		NegativePrompt string `json:"negative_prompt,omitempty"`
		Style          string `json:"style,omitempty"`
		Seed           uint64 `json:"seed,omitempty"`
		ImageUrl       string `json:"image_url,omitempty"`
	}
	TextToVideoResponse struct {
		Video Video `json:"video"`
	}
	Video struct {
		Url         string `json:"url"`
		ContentType string `json:"content_type,omitempty"`
		Filename    string `json:"filename,omitempty"`
		FileSize    int    `json:"file_size,omitempty"`
	}
)
