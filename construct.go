package oai_adapter

import (
	"github.com/jiu-u/oai-adapter/clients/base"
	"github.com/jiu-u/oai-adapter/clients/deepseek"
	"github.com/jiu-u/oai-adapter/clients/fal"
	"github.com/jiu-u/oai-adapter/clients/gemini_oai"
	"github.com/jiu-u/oai-adapter/clients/ollama_native"
	"github.com/jiu-u/oai-adapter/clients/ollama_oai"
	"github.com/jiu-u/oai-adapter/clients/openai"
	"github.com/jiu-u/oai-adapter/clients/siliconflow"
	"github.com/jiu-u/oai-adapter/clients/volcengine"
	"github.com/jiu-u/oai-adapter/clients/xai"
)

type AdapterConfig struct {
	AdapterType AdapterType
	ApiKey      string
	EndPoint    string
}

type AdapterType string

const (
	OpenAI      AdapterType = "OpenAI"
	DeepSeek    AdapterType = "DeepSeek"
	XAI         AdapterType = "XAI"
	SiliconFlow AdapterType = "SiliconFlow"

	Gemini     AdapterType = "Gemini"
	Gemini2OAI AdapterType = "Gemini2OAI"
	//GeminiNative AdapterType = "GeminiNative"

	Ollama       AdapterType = "Ollama"
	Ollama2OAI   AdapterType = "Ollama2OAI"
	OllamaNative AdapterType = "OllamaNative"

	VolcEngine AdapterType = "VolcEngine"

	Fal AdapterType = "Fal"
)

func NewAdapter(config *AdapterConfig) Adapter {
	switch config.AdapterType {
	case OpenAI:
		return openai.NewClient(config.EndPoint, config.ApiKey)
	case DeepSeek:
		return deepseek.NewClient(config.EndPoint, config.ApiKey)
	case SiliconFlow:
		return siliconflow.NewClient(config.EndPoint, config.ApiKey)
	case VolcEngine:
		return volcengine.NewClient(config.EndPoint, config.ApiKey)
	case Gemini, Gemini2OAI:
		return gemini_oai.NewClient(config.EndPoint, config.ApiKey)
	//case GeminiNative:
	//	return gemini_native.NewClient(config.EndPoint, config.ApiKey)
	case Ollama, Ollama2OAI:
		return ollama_oai.NewClient(config.EndPoint, config.ApiKey)
	case OllamaNative:
		return ollama_native.NewClient(config.EndPoint, config.ApiKey)
	case XAI:
		return xai.NewClient(config.EndPoint, config.ApiKey)

	case Fal:
		return fal.NewClient(config.EndPoint, config.ApiKey)
	default:
		return base.NewClient(config.EndPoint, config.ApiKey)
	}
}
