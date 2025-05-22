package api

type EmbeddingRequest struct {
	Input          any    `json:"input"`
	Model          string `json:"model"`
	User           string `json:"user,omitempty"`
	EncodingFormat string `json:"encoding_format,omitempty"`
	Dimensions     int    `json:"dimensions,omitempty"`
}

type EmbeddingResponse struct {
	Object string `json:"object"`
	Model  string `json:"model"`
	Data   []Data `json:"data"`
	Usage  Usage  `json:"usage"`
}

type Data struct {
	ID        string `json:"id,omitempty"`
	Object    string `json:"object"`
	Index     int    `json:"index"`
	Embedding []any  `json:"embedding,omitempty"`
}
