package v1

// https://docs.siliconflow.cn/cn/api-reference/rerank/create-rerank

type RerankRequest struct {
	Model           string   `json:"model" binding:"required"`
	Query           string   `json:"query" binding:"required"`
	Documents       []string `json:"documents"`
	TopN            int      `json:"top_n"`
	ReturnDocuments bool     `json:"return_documents,omitempty" binding:"omitempty,default=false"`
	MaxChunkPerDoc  int      `json:"max_chunk_per_doc,omitempty"`
	OverLapTokens   int      `json:"overlap_tokens,omitempty"`
}

type (
	RerankResponse struct {
		ID      string         `json:"id"`
		Results []RerankResult `json:"results"`
		Usage   *Usage         `json:"usage,omitempty"`
	}
	RerankResult struct {
		Document       RerankDocument `json:"document"`
		Index          int            `json:"index"`
		RelevanceScore float64        `json:"relevance_score"`
	}
	RerankDocument struct {
		Text string `json:"text"`
	}
)
