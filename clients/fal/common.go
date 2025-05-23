package fal

var Models = []string{
	"anthropic/claude-3.7-sonnet",
	"anthropic/claude-3.5-sonnet",
	"anthropic/claude-3-5-haiku",
	"anthropic/claude-3-haiku",
	"google/gemini-pro-1.5",
	"google/gemini-flash-1.5",
	"google/gemini-flash-1.5-8b",
	"google/gemini-2.0-flash-001",
	"meta-llama/llama-3.2-1b-instruct",
	"meta-llama/llama-3.2-3b-instruct",
	"meta-llama/llama-3.1-8b-instruct",
	"meta-llama/llama-3.1-70b-instruct",
	"openai/gpt-4o-mini",
	"openai/gpt-4o",
	"deepseek/deepseek-r1",
	"meta-llama/llama-4-maverick",
	"meta-llama/llama-4-scout",
}

const (
	HomeUrl           = "https://fal.run"
	PromptLimit       = 4800
	SystemPromptLimit = 4800
)
