package gemini

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	stdurl "net/url"
	"strings"
)

type Model struct {
	Name                       string   `json:"name"`
	Version                    string   `json:"version"`
	DisplayName                string   `json:"displayName"`
	Description                string   `json:"description"`
	InputTokenLimit            int      `json:"inputTokenLimit"`
	OutputTokenLimit           int      `json:"outputTokenLimit"`
	SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
	Temperature                float64  `json:"temperature,omitempty"`
	TopP                       float64  `json:"topP,omitempty"`
	TopK                       int      `json:"topK,omitempty"`
	MaxTemperature             float64  `json:"maxTemperature,omitempty"`
}

type Models struct {
	Models []Model `json:"models"`
}

func (c *Client) Models(ctx context.Context) ([]string, error) {
	baseURL := c.EndPoint + "/v1beta/models"
	params := stdurl.Values{}
	params.Set("pageSize", "1000")
	params.Set("key", c.APIKey)
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, _, err := c.DoJsonRequest(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	var respData Models
	err = sonic.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return nil, err
	}
	var models []string
	for _, model := range respData.Models {
		parts := strings.Split(model.Name, "/")
		if len(parts) <= 1 {
			fmt.Println("gemini 返回模型名称，字符串中没有 '/':", model.Name)
			continue
		}
		name := parts[len(parts)-1]
		models = append(models, name)
	}
	return models, nil
}
