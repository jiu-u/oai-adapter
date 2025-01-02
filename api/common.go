package api

type FunctionCall struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Params      any    `json:"parameters,omitempty"`
	Args        string `json:"arguments,omitempty"`
	Strict      bool   `json:"strict,omitempty"`
}

type Function = FunctionCall
type ToolCall struct {
	Id       string   `json:"id,omitempty"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}
