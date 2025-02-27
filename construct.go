package oai_adapter

import (
	"github.com/jiu-u/oai-adapter/clients/deepseek"
	"github.com/jiu-u/oai-adapter/clients/gemini"
	oainomodels "github.com/jiu-u/oai-adapter/clients/oai_no_models"
	"github.com/jiu-u/oai-adapter/clients/openai"
	"github.com/jiu-u/oai-adapter/clients/siliconflow"
	"github.com/jiu-u/oai-adapter/clients/siliconflow_free"
	stdurl "net/url"
)

type AdapterConfig struct {
	AdapterType  AdapterType
	ApiKey       string
	EndPoint     string
	ManualModels []string
	ProxyURL     *stdurl.URL
}

func NewAdapter(config *AdapterConfig) Adapter {
	switch config.AdapterType {
	case OpenAI:
		return openai.NewClient(config.EndPoint, config.ApiKey, config.ProxyURL)
	case Gemini:
		return gemini.NewClient(config.EndPoint, config.ApiKey, config.ProxyURL)
	case OAINoModels:
		return oainomodels.NewClient(config.EndPoint, config.ApiKey, config.ProxyURL, config.ManualModels)
	case SiliconFlow:
		return siliconflow.NewClient(config.EndPoint, config.ApiKey, config.ProxyURL)
	case SiliconFlowFree:
		return siliconflow_free.NewClient(config.EndPoint, config.ApiKey, config.ProxyURL)
	case DeepSeek:
		return deepseek.NewClient(config.EndPoint, config.ApiKey, config.ProxyURL)
	default:
		return openai.NewClient(config.EndPoint, config.ApiKey, config.ProxyURL)
	}
}
