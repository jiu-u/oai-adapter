package v1

import "encoding/json"

type (
	ResponsesResponse struct {
		CreatedAt          int                `json:"created_at,omitempty"`
		Error              any                `json:"error,omitempty"`
		ID                 string             `json:"id,omitempty"`
		IncompleteDetails  *IncompleteDetails `json:"incomplete_details,omitempty"`
		Instructions       string             `json:"instructions"`
		MaxOutputTokens    int                `json:"max_output_tokens,omitempty"`
		Metadata           json.RawMessage    `json:"metadata,omitempty"`
		Model              string             `json:"model,omitempty"`
		Object             string             `json:"object,omitempty"`
		Output             []ResponsesOutput  `json:"output,omitempty"`
		OutputText         string             `json:"output_text,omitempty"` // sdn only
		ParallelToolCalls  bool               `json:"parallel_tool_calls"`
		PreviousResponseID string             `json:"previous_response_id"`
		Reasoning          *Reasoning         `json:"reasoning,omitempty"`
		ServiceTier        string             `json:"service_tier,omitempty"`
		Status             string             `json:"status"`
		Temperature        float64            `json:"temperature"`
		Text               *ResponsesText     `json:"text,omitempty"`
		ToolChoice         json.RawMessage    `json:"tool_choice,omitempty"`
		Tools              []ResponsesTool    `json:"tools,omitempty"`
		TopP               float64            `json:"top_p"`
		Truncation         string             `json:"truncation"`
		Usage              *Usage             `json:"usage,omitempty"`
		User               string             `json:"user,omitempty"`
	}
	ResponsesStreamResponse struct {
		Type     string             `json:"type"`
		Response *ResponsesResponse `json:"response"`
	}
)

type (
	IncompleteDetails struct {
		Reasoning string `json:"reasoning"`
	}
	ResponsesOutput struct {
		ID      string          `json:"id"`
		Type    string          `json:"type,omitempty"`
		Content json.RawMessage `json:"content"` // string or array(InputContent)
		Role    string          `json:"role,omitempty"`
		Status  string          `json:"status,omitempty"`
		Queries any             `json:"queries,omitempty"`
		Results any             `json:"results,omitempty"`
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
)
