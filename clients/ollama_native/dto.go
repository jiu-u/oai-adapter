package ollama_native

import (
	v1 "github.com/jiu-u/oai-adapter/api/v1"
	"time"
)

type (
	OllamaCompletionRequest struct {
		Model     string   `json:"model" binding:"required"`
		Prompt    string   `json:"prompt"`
		Stream    bool     `json:"stream,omitempty"`
		Suffix    string   `json:"suffix,omitempty"`
		Images    []string `json:"images,omitempty"`
		Format    any      `json:"format,omitempty"`
		Options   *Options `json:"options,omitempty"`
		System    string   `json:"system,omitempty"`
		Template  string   `json:"template,omitempty"`
		Raw       bool     `json:"raw,omitempty"`
		KeepAlive string   `json:"keep_alive,omitempty"`
	}
	Options struct {
		// NumCtx 设置用于生成下一个标记的上下文窗口的大小。（默认：2048）
		NumCtx int64 `json:"num_ctx,omitempty"`
		// RepeatLastN 设置模型回溯的长度以防止重复。（默认：64，0 = 禁用，-1 = num_ctx）
		RepeatLastN int `json:"repeat_last_n,omitempty"`
		// RepeatPenalty 设置对重复内容的惩罚强度。较高的值（例如 1.5）会更强地惩罚重复内容，而较低的值（例如 0.9）则更宽容。（默认值：1.1）
		RepeatPenalty float64 `json:"repeat_penalty,omitempty"`
		// Temperature 模型的温度。提高温度会使模型的回答更具创造性。（默认值：0.8）
		Temperature float64 `json:"temperature,omitempty"` // Temperature 控制输出的随机性和创造性
		// Seed 设置用于生成的随机数种子。将此设置为特定数字将使模型对相同的提示生成相同的文本。（默认：0）
		Seed int64 `json:"seed,omitempty"`
		// Stop 设置用于终止生成的字符。如果提供了 stop，则在每个完成的 token 中，将停止发生在 stop 字符串出现的位置。（默认：空字符串）
		Stop any `json:"stop,omitempty"`
		// NumPredict 生成文本时预测的最大令牌数。（默认：-1，无限生成）
		NumPredict int `json:"num_predict,omitempty"`
		// TopK 降低生成无意义内容的概率。较高的值（例如 100）会给出更多样化的答案，而较低的值（例如 10）则更为保守。（默认值：40）
		TopK int `json:"top_k,omitempty"`
		// TopP 降低生成无意义内容的概率。较高的值（例如 0.8）会给出更多样化的答案，而较低的值（例如 0.2）则更为保守。（默认值：0.9）
		TopP float64 `json:"top_p,omitempty"`
		// MinP top_p 的替代方案，旨在确保质量和多样性之间的平衡。参数 p 表示一个标记被考虑的最小概率，相对于最可能标记的概率。例如，当 p=0.05 且最可能标记的概率为 0.9 时，值小于 0.045 的 logits 会被过滤掉。（默认值：0.0）
		MinP             float64 `json:"min_p,omitempty"`
		NumKeep          int     `json:"num_keep,omitempty"`
		TypicalP         float64 `json:"typical_p,omitempty"`
		PresencePenalty  float64 `json:"presence_penalty,omitempty"`
		FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
		PenalizeNewline  bool    `json:"penalize_newline,omitempty"`
		Numa             bool    `json:"numa,omitempty"`
		NumBatch         int     `json:"num_batch,omitempty"`
		NumGpu           int     `json:"num_gpu,omitempty"`
		MianGpu          int     `json:"mian_gpu,omitempty"`
		UseMMap          bool    `json:"use_mmap,omitempty"`
		NumThread        int     `json:"num_thread,omitempty"`

		//Echo             bool    `json:"echo,omitempty"`
		//Logprobs         any     `json:"logprobs,omitempty"`
		//MaxTokens        int64   `json:"max_tokens,omitempty"`
		//BestOf           int64   `json:"best_of,omitempty"`
		//LogitBias        any     `json:"logit_bias,omitempty"`
		//N                any     `json:"n,omitempty"`
		//User             string  `json:"user,omitempty"`
	}
	OllamaCompletionResponse struct {
		Model              string    `json:"model"`
		CreatedAt          time.Time `json:"created_at"`
		Response           string    `json:"response"`
		Done               bool      `json:"done"`
		DoneReason         string    `json:"done_reason，omitempty"`
		Context            []int     `json:"context"`
		TotalDuration      int64     `json:"total_duration"`
		LoadDuration       int64     `json:"load_duration"`
		PromptEvalCount    int       `json:"prompt_eval_count"`
		PromptEvalDuration int64     `json:"prompt_eval_duration"`
		EvalCount          int       `json:"eval_count"`
		EvalDuration       int64     `json:"eval_duration"`
	}
	OllamaCompletionStreamResponse = OllamaCompletionResponse
)

type (
	OllamaChatCompletionRequest struct {
		Model     string          `json:"model" binding:"required"`
		Messages  []OllamaMessage `json:"messages,omitempty"`
		Tools     []v1.Tool       `json:"tools,omitempty"`
		Format    any             `json:"format,omitempty"`
		Options   *Options        `json:"options,omitempty"`
		Stream    bool            `json:"stream,omitempty"`
		KeepAlive string          `json:"keep_alive,omitempty"`
	}
	OllamaMessage struct {
		Role      string        `json:"role" binding:"required"`
		Content   string        `json:"content" binding:"required"`
		Images    []string      `json:"images,omitempty"`
		ToolCalls []v1.ToolCall `json:"tool_calls,omitempty"`
	}
	OllamaChatCompletionResponse struct {
		Model              string          `json:"model"`
		CreatedAt          time.Time       `json:"created_at"`
		Message            []OllamaMessage `json:"message"`
		Done               bool            `json:"done"`
		DoneReason         string          `json:"done_reason，omitempty"`
		Context            []int           `json:"context,omitempty"`
		TotalDuration      int64           `json:"total_duration"`
		LoadDuration       int64           `json:"load_duration"`
		PromptEvalCount    int             `json:"prompt_eval_count"`
		PromptEvalDuration int64           `json:"prompt_eval_duration"`
		EvalCount          int             `json:"eval_count"`
		EvalDuration       int64           `json:"eval_duration"`
	}
)

type (
	OllamaEmbeddingRequest struct {
		Model     string   `json:"model"`
		Input     any      `json:"input"`
		Truncate  bool     `json:"truncate,omitempty"`
		Options   *Options `json:"options,omitempty"`
		KeepAlive string   `json:"keep_alive,omitempty"`
	}
	OllamaEmbeddingResponse struct {
		Model           string  `json:"model"`
		Embeddings      [][]any `json:"embeddings"`
		TotalDuration   int64   `json:"total_duration"`
		LoadDuration    int64   `json:"load_duration"`
		PromptEvalCount int     `json:"prompt_eval_count"`
	}
)
