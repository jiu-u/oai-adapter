package v1

import "mime/multipart"

type (
	ImageGenerateRequest struct {
		Prompt            string `json:"prompt" binding:"required"`
		Background        string `json:"background,omitempty"`
		Model             string `json:"model" binding:"required"`
		Moderation        string `json:"moderation,omitempty"`
		N                 int    `json:"n,omitempty"`
		OutputCompression int    `json:"output_compression,omitempty"`
		OutputFormat      string `json:"output_format,omitempty"`
		Quality           string `json:"quality,omitempty"`
		ResponseFormat    string `json:"response_format,omitempty"`
		Size              string `json:"size,omitempty"`
		Style             string `json:"style,omitempty"` //  one of vivid or natural
		User              string `json:"user,omitempty"`
	}
	ImageGenerateResponse struct {
		Created int64          `json:"created"`
		Data    []ImageGenData `json:"data"`
	}
	ImageUsage struct {
		InputTokens       int                     `json:"input_tokens"`
		InputTokenDetails *ImageInputTokenDetails `json:"input_token_details,omitempty"`
		OutputTokens      int                     `json:"output_tokens"`
		TotalTokens       int                     `json:"total_tokens"`
	}
	ImageInputTokenDetails struct {
		ImageTokens int `json:"image_tokens"`
		TextTokens  int `json:"text_tokens"`
	}
	ImageGenData struct {
		URL           string `json:"url,omitempty"`
		B64JSON       string `json:"b64_json,omitempty"`
		RevisedPrompt string `json:"revised_prompt,omitempty"`
	}
)

type (
	ImageEditRequest struct {
		Image          []*multipart.FileHeader `form:"image" binding:"required"`
		Prompt         string                  `form:"prompt" binding:"required"`
		Background     string                  `form:"background,omitempty"`
		Mask           *multipart.FileHeader   `form:"mask,omitempty"`
		Model          string                  `form:"model" binding:"required"`
		N              int                     `form:"n,omitempty"`
		Quality        string                  `form:"quality,omitempty"`
		ResponseFormat string                  `form:"response_format,omitempty"`
		Size           string                  `form:"size,omitempty"`
		User           string                  `form:"user,omitempty"`
	}
	ImageEditResponse = ImageGenerateResponse
)

type (
	ImageVariationRequest struct {
		Image          *multipart.FileHeader `form:"image" binding:"required"`
		Model          string                `form:"model" binding:"required"`
		N              int                   `form:"n,omitempty"`
		Size           string                `form:"size,omitempty"`
		ResponseFormat string                `form:"response_format,omitempty"`
		User           string                `form:"user,omitempty"`
	}
	ImageVariationResponse = ImageGenerateResponse
)
