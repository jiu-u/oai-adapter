package v1

import (
	"encoding/json"
	"fmt"
)

type (
	//Include json.RawMessage `json:"include,omitempty"`

	ResponsesRequest struct {
		Input              json.RawMessage `json:"input"` // string、[]ResponsesInput
		Model              string          `json:"model" binding:"required"`
		Include            []string        `json:"include,omitempty"`
		Instructions       string          `json:"instructions,omitempty"`
		MaxOutputTokens    uint            `json:"max_output_tokens,omitempty"`
		Metadata           any             `json:"metadata,omitempty"`
		ParallelToolCalls  bool            `json:"parallel_tool_calls,omitempty"`
		PreviousResponseId string          `json:"previous_response_id,omitempty"`
		Reasoning          *Reasoning      `json:"reasoning,omitempty"`
		ServiceTier        string          `json:"service_tier,omitempty"`
		Store              bool            `json:"store,omitempty"`
		Steam              bool            `json:"steam,omitempty"`
		Temperature        float64         `json:"temperature,omitempty"`
		Text               *ResponsesText  `json:"text,omitempty"`
		ToolChoice         json.RawMessage `json:"tool_choice,omitempty"` // string or object(ToolChoice)
		Tools              []ResponsesTool `json:"tools,omitempty"`
		TopP               float64         `json:"top_p,omitempty"`
		Truncation         string          `json:"truncation,omitempty"`
		User               string          `json:"user,omitempty"`
	}
	ResponsesText struct {
		Type string `json:"type"`
		//json schema
		Name        string `json:"name,omitempty"`
		Schema      any    `json:"schema,omitempty"`
		Description string `json:"description,omitempty"`
		Strict      bool   `json:"strict,omitempty"`
	}
	ResponsesTool struct {
		Type string `json:"type"`
		// Web Search
		SearchContextSize string `json:"search_context_size,omitempty"`
		UserLocation      any    `json:"user_location,omitempty"`
		// File Search
		VectorStoreIds []string `json:"vector_store_ids,omitempty"`
		Filters        any      `json:"filters,omitempty"`
		MaxNumResults  uint     `json:"max_num_results,omitempty"`
		RankingOptions any      `json:"ranking_options,omitempty"`
		// Computer Use
		DisplayHeight uint   `json:"display_height,omitempty"`
		DisplayWidth  uint   `json:"display_width,omitempty"`
		Environment   string `json:"environment,omitempty"`
		// Function
		Name        string `json:"name,omitempty"`
		Parameters  any    `json:"parameters,omitempty"`
		Strict      bool   `json:"strict,omitempty"`
		Description string `json:"description,omitempty"`
	}
)

const (
	MessageInputType            = "message"
	TextInputType               = "text"
	ItemReferenceInputType      = "item_reference"
	FileSearchCallInputType     = "file_search_call"
	ComputerCallInputType       = "computer_call"
	ComputerCallOutputInputType = "computer_call_output"
	WebSearchCallInputType      = "web_search_call"
	FunctionCallInputType       = "function_call"
	FunctionCallOutputInputType = "function_call_output"
	ReasoningInputType          = "reasoning"
)

func (r *ResponsesRequest) IsStringInput() bool {
	var stringInput string
	if err := json.Unmarshal(r.Input, &stringInput); err == nil {
		return true
	}
	return false
}

func (r *ResponsesRequest) StringInput() string {
	var stringInput string
	if err := json.Unmarshal(r.Input, &stringInput); err == nil {
		return stringInput
	}
	return string(r.Input)
}

func (r *ResponsesRequest) ParseInput() ([]ResponsesInput, error) {
	var err error
	var inputList []ResponsesInput
	if r.IsStringInput() {
		inputList = append(inputList, ResponsesInput{
			ID:      "input",
			Type:    MessageInputType,
			Content: json.RawMessage(r.StringInput()),
			Role:    "user",
		})
		return inputList, nil
	}
	if err := json.Unmarshal(r.Input, &inputList); err == nil {
		return inputList, nil
	}

	return nil, fmt.Errorf("parse input error: %w", err)
}

type ResponsesInput struct {
	ID      string          `json:"id"`
	Type    string          `json:"type,omitempty"`
	Content json.RawMessage `json:"content"` // string or array(InputContent)
	Role    string          `json:"role,omitempty"`
	Status  string          `json:"status,omitempty"`
	// File Search Tool
	Queries any `json:"queries,omitempty"`
	Results any `json:"results,omitempty"`
	// Computer Tool call
	Action              any `json:"action,omitempty"`
	CallId              any `json:"call_id,omitempty"`
	PendingSafetyChecks any `json:"pending_safety_checks,omitempty"`
	// function tool call
	Arguments string `json:"arguments,omitempty"`
	Name      string `json:"name,omitempty"`
	// function tool call output
	Output string `json:"output,omitempty"`
	// reasoning
	Summary          any `json:"summary,omitempty"`
	EncryptedContent any `json:"encrypted_content,omitempty"`
}

func (r *ResponsesInput) IsStringContent() bool {
	var stringContent string
	if err := json.Unmarshal(r.Content, &stringContent); err == nil {
		return true
	}
	return false
}

func (r *ResponsesInput) StringContent() string {
	var stringContent string
	if err := json.Unmarshal(r.Content, &stringContent); err == nil {
		return stringContent
	}
	return string(r.Content)
}

func (r *ResponsesInput) ParseContent() ([]InputContent, error) {
	var err error
	var contentList []InputContent
	if r.IsStringContent() {
		contentList = append(contentList, InputContent{
			Type: InputContentTypeText,
			Text: r.StringContent(),
		})
		return contentList, nil
	}
	var inputContentList []InputContent
	if err := json.Unmarshal(r.Content, &inputContentList); err == nil {
		return inputContentList, nil
	}
	return nil, fmt.Errorf("parse content error: %w", err)
}

type InputContent struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
	// input image
	Detail   string `json:"detail,omitempty"`
	FileId   string `json:"file_id,omitempty"`
	ImageUrl string `json:"image_url,omitempty"`
	// input file
	FileData string `json:"file_data,omitempty"`
	Filename string `json:"filename,omitempty"`
	// output message content output text
	Annotations any `json:"annotations,omitempty"`
	// output message content refusal
	Refusal any `json:"refusal,omitempty"`
}

const (
	InputContentTypeText       = "input_text"
	InputContentTypeImage      = "input_image"
	InputContentTypeInputAudio = "input_file"
	InputContentTypeFile       = "file"
	InputContentTypeOutputText = "output_text"
	InputContentTypeRefusal    = "refusal"
)
