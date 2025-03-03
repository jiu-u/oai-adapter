package oai_adapter

type AdapterType string

const (
	OpenAI          AdapterType = "openai"
	Gemini          AdapterType = "gemini"
	OAINoModels     AdapterType = "oaiNoModels"
	SiliconFlow     AdapterType = "siliconflow"
	SiliconFlowFree AdapterType = "siliconflowFree"
	OLLAMA          AdapterType = "ollama"
	DeepSeek        AdapterType = "deepseek"
	XAI             AdapterType = "xai"
)
