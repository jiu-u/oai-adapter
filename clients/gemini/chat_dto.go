package gemini

import "encoding/json"

type (
	ChatRequest struct {
		Contents          []Content         `json:"contents"`                    //必需。与模型的当前对话内容。
		Tools             []Tool            `json:"tools,omitempty"`             //可选。Model 可能用于生成下一个响应的 Tools 列表。
		ToolConfig        *ToolConfig       `json:"toolConfig,omitempty"`        //可选。请求中指定的任何 Tool 的工具配置
		SafetySettings    []SafetySettings  `json:"safetySettings,omitempty"`    //可选。用于屏蔽不安全内容的唯一 SafetySetting 实例列表。
		SystemInstruction *Content          `json:"systemInstruction,omitempty"` //可选。开发者设置的系统说明。目前仅支持文字广告。
		GenerationConfig  *GenerationConfig `json:"generationConfig,omitempty"`  //可选。模型生成和输出的配置选项。
		CachedContent     string            `json:"cachedContent,omitempty"`     //可选。已缓存的内容的名称，用作上下文以提供预测。
	}
	Content struct {
		Role  string `json:"role"`
		Parts []Part `json:"parts"`
	}
	Part struct {
		Text                string               `json:"text,omitempty"`
		InlineData          *Blob                `json:"inlineData,omitempty"`
		FunctionCall        *FunctionCall        `json:"functionCall,omitempty"`
		FunctionResponse    *FunctionResponse    `json:"functionResponse,omitempty"`
		FileData            *FileData            `json:"fileData,omitempty"`
		ExecutableCode      *ExecutableCode      `json:"executableCode,omitempty"`
		CodeExecutionResult *CodeExecutionResult `json:"codeExecutionResult,omitempty"`
	}
	Blob struct {
		MimeType string `json:"mimeType"` //来源数据的 IANA 标准 MIME 类型。示例：- image/png - image/jpeg 如果提供的 MIME 类型不受支持，系统会返回错误。
		Data     string `json:"data"`     //媒体格式的原始字节。 使用 base64 编码的字符串。
	}
	FunctionCall struct {
		Name string          `json:"name"`           // FunctionDeclaration.name
		Args json.RawMessage `json:"args,omitempty"` // 参数及其值，使用 json.RawMessage 可以存储任意 JSON 结构
	}
	FunctionResponse struct {
		Name     string          `json:"name"`               // FunctionDeclaration.name
		Response json.RawMessage `json:"response,omitempty"` // 函数的输出，使用 json.RawMessage 可以存储任意 JSON 结构
	}
	FileData struct {
		MimeType *string `json:"mimeType,omitempty"`
		FileURL  string  `json:"fileUri"`
	}
	// ExecutableCode 表示由模型生成的要执行的代码。
	ExecutableCode struct {
		Code     string `json:"code"`
		Language string `json:"language"`
	}
	// CodeExecutionResult 表示执行 ExecutableCode 的结果。
	CodeExecutionResult struct {
		OutCome string  `json:"outcome"`
		Output  *string `json:"output,omitempty"`
	}
	Tool struct {
	}
	ToolConfig     struct{}
	safetySettings struct{}
	SafetySettings struct {
	}

	GenerationConfig struct {
		StopSequences    []string `json:"stopSequences,omitempty"`
		ResponseMimeType string   `json:"responseMimeType,omitempty"`
		ResponseSchema   *Schema  `json:"responseSchema,omitempty"`
		CandidateCount   int      `json:"candidateCount,omitempty"`
		MaxOutputTokens  int      `json:"maxOutputTokens,omitempty"`
		Temperature      float64  `json:"temperature,omitempty"`
		TopP             float64  `json:"topP,omitempty"`
		TopK             int      `json:"topK,omitempty"`
		PresencePenalty  float64  `json:"presencePenalty,omitempty"`
		FrequencyPenalty float64  `json:"frequencyPenalty,omitempty"`
		ResponseLogprobs bool     `json:"responseLogprobs,omitempty"`
		Logprobs         int      `json:"logprobs,omitempty"`
	}
	Schema struct {
		Type        string         `json:"type"` // 必需
		Format      *string        `json:"format,omitempty"`
		Description *string        `json:"description,omitempty"`
		Nullable    *bool          `json:"nullable,omitempty"`
		Enum        []string       `json:"enum,omitempty"`
		MaxItems    *int64         `json:"maxItems,omitempty"`
		MinItems    *int64         `json:"minItems,omitempty"`
		Properties  map[string]any `json:"properties,omitempty"`
		Required    []string       `json:"required,omitempty"`
		Items       *Schema        `json:"items,omitempty"`
	}
)

type (
	GenerateContentResponse struct {
		Candidates    []Candidate   `json:"candidates"`
		UsageMetadata UsageMetadata `json:"usageMetadata"`
		ModelVersion  string        `json:"modelVersion"`
		//PromptFeedback PromptFeedback `json:"promptFeedback"`
	}

	Candidate struct {
		Content      Content `json:"content"`
		FinishReason string  `json:"finishReason"`
		TokenCount   int     `json:"tokenCount"`
		AvgLogprobs  float64 `json:"avgLogprobs"`
		Index        int     `json:"index"`
		//SafetyRatings         []SafetyRating         `json:"safetyRatings"`
		//CitationMetadata      CitationMetadata       `json:"citationMetadata"`
		//GroundingAttributions []GroundingAttribution `json:"groundingAttributions"`
		//GroundingMetadata     GroundingMetadata      `json:"groundingMetadata"`
		//LogprobsResult LogprobsResult `json:"logprobsResult"`
	}
	// PromptFeedback 是提示的反馈元数据。
	PromptFeedback struct {
		BlockReason   string         `json:"blockReason"`
		SafetyRatings []SafetyRating `json:"safetyRatings"`
	}
	// UsageMetadata 是关于生成请求的令牌使用情况的元数据。
	UsageMetadata struct {
		PromptTokenCount        int `json:"promptTokenCount"`
		CachedContentTokenCount int `json:"cachedContentTokenCount"`
		CandidatesTokenCount    int `json:"candidatesTokenCount"`
		TotalTokenCount         int `json:"totalTokenCount"`
	}
	SafetyRating struct {
		Category    string `json:"category"`
		Probability string `json:"probability"`
		Blocked     bool   `json:"blocked"`
	}
	// CitationMetadata 是内容的引用信息。
	CitationMetadata struct {
		CitationSources []CitationSource `json:"citationSources"`
	}

	// CitationSource 是对特定回答的一部分的引用。
	CitationSource struct {
		StartIndex int    `json:"startIndex"`
		EndIndex   int    `json:"endIndex"`
		Uri        string `json:"uri"`
		License    string `json:"license"`
	}
	// GroundingAttribution 是对回答做出贡献的来源的归因。
	GroundingAttribution struct {
		SourceId AttributionSourceId `json:"sourceId"`
		Content  Content             `json:"content"`
	}

	// AttributionSourceId 是对此归因做出贡献的来源的标识符。
	AttributionSourceId struct {
		GroundingPassage       *GroundingPassageId     `json:"groundingPassage,omitempty"`
		SemanticRetrieverChunk *SemanticRetrieverChunk `json:"semanticRetrieverChunk,omitempty"`
	}
	// GroundingPassageId 是 GroundingPassage 中某个部分的标识符。
	GroundingPassageId struct {
		PassageId string `json:"passageId"`
		PartIndex int    `json:"partIndex"`
	}

	// SemanticRetrieverChunk 是通过语义检索器检索的 Chunk 的标识符。
	SemanticRetrieverChunk struct {
		Source string `json:"source"`
		Chunk  string `json:"chunk"`
	}
	// GroundingMetadata 是启用接地时返回的元数据。
	GroundingMetadata struct {
		GroundingChunks   []GroundingChunk   `json:"groundingChunks"`
		GroundingSupports []GroundingSupport `json:"groundingSupports"`
		WebSearchQueries  []string           `json:"webSearchQueries"`
		SearchEntryPoint  SearchEntryPoint   `json:"searchEntryPoint"`
		RetrievalMetadata RetrievalMetadata  `json:"retrievalMetadata"`
	}

	// GroundingChunk 是接地块。
	GroundingChunk struct {
		Web *Web `json:"web,omitempty"`
	}

	// Web 是来自网络的文件块。
	Web struct {
		Uri   string `json:"uri"`
		Title string `json:"title"`
	}

	// GroundingSupport 是接地支持。
	GroundingSupport struct {
		GroundingChunkIndices []int     `json:"groundingChunkIndices"`
		ConfidenceScores      []float64 `json:"confidenceScores"`
		Segment               Segment   `json:"segment"`
	}

	// Segment 是内容的片段。
	Segment struct {
		PartIndex  int    `json:"partIndex"`
		StartIndex int    `json:"startIndex"`
		EndIndex   int    `json:"endIndex"`
		Text       string `json:"text"`
	}

	// SearchEntryPoint 是 Google 搜索入口点。
	SearchEntryPoint struct {
		RenderedContent string `json:"renderedContent"`
		SdkBlob         string `json:"sdkBlob"`
	}

	// RetrievalMetadata 是与基准流程中检索相关的元数据。
	RetrievalMetadata struct {
		GoogleSearchDynamicRetrievalScore float64 `json:"googleSearchDynamicRetrievalScore"`
	}
	// LogprobsResult 是 logprobs 结果。
	LogprobsResult struct {
		TopCandidates    []TopCandidates `json:"topCandidates"`
		ChosenCandidates []Candidate2    `json:"chosenCandidates"`
	}

	// TopCandidates 是在每个解码步骤中具有最高对数概率的候选网络。
	TopCandidates struct {
		Candidates []Candidate2 `json:"candidates"`
	}

	// Candidate2 是 logprobs 词元和得分的候选对象。
	Candidate2 struct {
		Token          string  `json:"token"`
		TokenId        int     `json:"tokenId"`
		LogProbability float64 `json:"logProbability"`
	}
)
