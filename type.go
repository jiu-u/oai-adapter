package oai_adapter

type AdapterType string

const (
	OpenAI          AdapterType = "openai"
	Gemini          AdapterType = "gemini"
	OAINoModels     AdapterType = "oaiNoModels"
	SiliconFlow     AdapterType = "siliconflow"
	SiliconFlowFree AdapterType = "siliconflowFree"
)
