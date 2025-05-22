package v1

type (
	EmbeddingsRequest struct {
		Input          any    `json:"input"`
		Model          string `json:"model"`
		Dimensions     int    `json:"dimensions,omitempty"`
		EncodingFormat string `json:"encoding_format,omitempty" binding:"omitempty,oneof=base64 float,default=float"`
		User           string `json:"user,omitempty"`
	}
)

type (
	EmbeddingsResponse struct {
		Object string           `json:"object"`
		Model  string           `json:"model"`
		Data   []EmbeddingsData `json:"data"`
		Usage  Usage            `json:"usage"`
	}
	EmbeddingsData struct {
		ID        string `json:"id,omitempty"`
		Object    string `json:"object"`
		Index     int    `json:"index"`
		Embedding []any  `json:"embedding,omitempty"`
	}
)
