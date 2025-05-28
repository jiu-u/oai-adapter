package v1

import "mime/multipart"

type (
	AudioSpeechRequest struct {
		Input          string  `json:"input" binding:"required"`
		Model          string  `json:"model" binding:"required" default:"tts-1"`
		Voice          string  `json:"voice" binding:"required" default:"alloy"`
		ResponseFormat string  `json:"response_format,omitempty"` // map、opus、acc、flac、wav、pcm
		Speed          float64 `json:"speed,omitempty"`
		Instructions   string  `json:"instructions,omitempty"`
	}
	// 二进制文件
	//AudioSpeechResponse struct{}
)

type (
	TranscriptionRequest struct {
		File                   *multipart.FileHeader `json:"file,omitempty" form:"file" binding:"required"`
		Model                  string                `json:"model" form:"model" binding:"required"`
		ChunkingStrategy       string                `json:"chunking_strategy,omitempty" form:"chunking_strategy,omitempty"`
		Include                []string              `form:"include,omitempty"`
		Language               string                `form:"language,omitempty"`
		Prompt                 string                `form:"prompt,omitempty"`
		ResponseFormat         string                `form:"response_format,omitempty"`
		Temperature            float64               `form:"temperature,omitempty"` // 温度
		TimestampGranularities []string              `form:"timestamp_granularities,omitempty"`
	}
	TranscriptionLogprobs struct {
		Token   string  `json:"token"`
		Logprob float64 `json:"logprob"`
		Bytes   []int   `json:"bytes"`
	}
	TranscriptionResponse struct {
		Logprobs []TranscriptionLogprobs `json:"logprobs,omitempty"`
		Text     string                  `json:"text"`
	}
	TranscriptionWord struct {
		Word  string  `json:"word"`
		Start float64 `json:"start"`
		End   float64 `json:"end"`
	}
	TranscriptionSegment struct {
		AverageLogprob   float64 `json:"average_logprob"`
		CompressionRatio float64 `json:"compression_ratio"`
		End              float64 `json:"end"`
		Id               int     `json:"id"`
		NoSpeechProb     float64 `json:"no_speech_prob"`
		Seek             int     `json:"seek"`
		Start            float64 `json:"start"`
		Temperature      float64 `json:"temperature"`
		Text             string  `json:"text"`
		Tokens           []int   `json:"tokens,omitempty"`
	}
	TranscriptionVerboseResponse struct {
		Duration float64                `json:"duration"`
		Language string                 `json:"language"`
		Text     string                 `json:"text"`
		Segments []TranscriptionSegment `json:"segments,omitempty"`
		Words    []TranscriptionWord    `json:"words,omitempty"`
	}
	TranscriptionStreamResponse struct {
		Delta    string                  `json:"delta"`
		Type     string                  `json:"type"`
		Logprobs []TranscriptionLogprobs `json:"logprobs,omitempty"`
	}
)

type (
	TranslationRequest struct {
		File           *multipart.FileHeader `form:"file" binding:"required"`
		Model          string                `form:"model" binding:"required"`
		Prompt         string                `form:"prompt,omitempty"`
		ResponseFormat string                `form:"response_format,omitempty"`
		Temperature    float64               `form:"temperature,omitempty"` // 温度
	}
	// 二进制文件
	//AudioSpeechResponse struct{}
)
