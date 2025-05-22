package gemini

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/jiu-u/oai-adapter/api"
	"github.com/jiu-u/oai-adapter/tools"
	"io"
	"net/http"
	stdurl "net/url"
	"strconv"
	"strings"
	"time"
)

var toGeminiRoleMap = map[string]string{
	"system":    "user",
	"user":      "user",
	"assistant": "model",
	"developer": "user",
	"model":     "assistant",
}

var toOpenAiRoleMap = map[string]string{
	"user":      "user",
	"model":     "assistant",
	"developer": "developer",
	"system":    "system",
}

func message2Parts(msg *api.Message) ([]Part, error) {
	parts := make([]Part, 0)
	part := new(Part)
	if msg.IsStringContent() {
		part.Text = msg.StringContent()
		parts = append(parts, *part)
		return parts, nil
	}

	mediaContents, err := msg.ParseContent()
	if err != nil {
		return nil, err
	}
	for _, mediaContent := range mediaContents {
		part = new(Part)
		switch mediaContent.Type {
		case api.ContentTypeText:
			part.Text = mediaContent.Text
			parts = append(parts, *part)
			continue
		case api.ContentTypeImageURL:
			if mediaContent.ImageUrl.Url == "" {
				return nil, errors.New("image_url is empty")
			}
			f, err := tools.NewImageFileData(mediaContent.ImageUrl.Url, true)
			if err != nil {
				return nil, err
			}
			b64Parts := strings.SplitN(f.URL, ",", 2)
			part.InlineData = &Blob{
				MimeType: f.MIMEType,
				Data:     b64Parts[1],
			}
			parts = append(parts, *part)
			continue
		case api.ContentTypeInputAudio:
			if mediaContent.InputAudio.Data == "" {
				return nil, errors.New("audio_data is empty")
			}
			if mediaContent.InputAudio.Format == "" {
				return nil, errors.New("audio_format is empty")
			}
			part.InlineData = &Blob{
				MimeType: "audio/" + mediaContent.InputAudio.Format,
				Data:     mediaContent.InputAudio.Data,
			}
			parts = append(parts, *part)
			continue
		default:
			return nil, errors.New("unknown content type")
		}
	}
	return parts, err
}

func messagesToContents(messages []api.Message) ([]Content, error) {
	contents := make([]Content, 0, len(messages))
	for _, msg := range messages {
		parts, err := message2Parts(&msg)
		//parts, err := content2Part(msg.Content)
		if err != nil {
			return nil, err
		}
		content := Content{
			Role:  toGeminiRoleMap[msg.Role],
			Parts: parts,
		}
		contents = append(contents, content)
	}
	return contents, nil
}

func (c *Client) ConvertChatRequest(req *api.ChatRequest) (*ChatRequest, error) {
	resp := &ChatRequest{
		GenerationConfig: &GenerationConfig{
			Temperature:      req.Temperature,
			TopP:             req.TopP,
			MaxOutputTokens:  req.MaxTokens,
			PresencePenalty:  req.PresencePenalty,
			FrequencyPenalty: req.FrequencyPenalty,
		},
	}
	contents, err := messagesToContents(req.Messages)
	if err != nil {
		return nil, err
	}
	resp.Contents = contents
	return resp, nil
}

func (c *Client) ConvertChatResponse(resp *GenerateContentResponse) (*api.ChatCompletionNoStreamResponse, error) {
	now := time.Now().Unix()
	newResp := &api.ChatCompletionNoStreamResponse{
		ID:                "chatcmpl-" + strconv.FormatInt(now, 10),
		Object:            "chat.completion",
		Created:           now,
		Model:             resp.ModelVersion,
		SystemFingerprint: "fp_" + strconv.FormatInt(now, 10),
		Choices:           nil,
		Usage: api.Usage{
			PromptTokens:     resp.UsageMetadata.PromptTokenCount,
			CompletionTokens: resp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      resp.UsageMetadata.TotalTokenCount,
		},
	}
	choices := make([]api.Choice, 0, len(resp.Candidates))
	index := 0
	for _, candidate := range resp.Candidates {
		role := candidate.Content.Role
		role = toOpenAiRoleMap[role]
		for _, part := range candidate.Content.Parts {
			message := new(api.Message)
			message.Role = role
			if part.Text == "" {
				continue
			}
			str, err := json.Marshal(part.Text)
			if err != nil {
				return nil, err
			}
			message.Content = str
			choices = append(choices, api.Choice{
				Index:        index,
				FinishReason: candidate.FinishReason,
				Message:      *message,
			})
			index++
		}

	}
	newResp.Choices = choices
	return newResp, nil
}

func (c *Client) ConvertStreamChatResponse(resp *GenerateContentResponse) (*api.ChatCompletionStreamResponse, error) {
	now := time.Now().Unix()
	newResp := &api.ChatCompletionStreamResponse{
		ID:      "chatcmpl-" + strconv.FormatInt(now, 10),
		Object:  "chat.completion.chunk",
		Created: now,
		Model:   resp.ModelVersion,
		Choices: nil,
		Usage: api.Usage{
			PromptTokens:     resp.UsageMetadata.PromptTokenCount,
			CompletionTokens: resp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      resp.UsageMetadata.TotalTokenCount,
		},
	}
	choices := make([]api.ChoiceWithDelta, 0, len(resp.Candidates))
	index := 0
	for _, candidate := range resp.Candidates {
		role := candidate.Content.Role
		role = toOpenAiRoleMap[role]
		for _, part := range candidate.Content.Parts {
			delta := new(api.Delta)
			delta.Role = role
			if part.Text == "" {
				continue
			}
			delta.Content = part.Text
			choices = append(choices, api.ChoiceWithDelta{
				Index:        index,
				FinishReason: candidate.FinishReason,
				Delta:        *delta,
			})
			index++
		}
	}
	newResp.Choices = choices
	return newResp, nil
}

func (c *Client) convertNoStreamChatCompletions(resp io.ReadCloser, header http.Header) (io.ReadCloser, http.Header, error) {
	bodyBytes, err := io.ReadAll(resp)
	if err != nil {
		return nil, nil, err
	}
	var respBody GenerateContentResponse
	err = sonic.Unmarshal(bodyBytes, &respBody)
	if err != nil {
		fmt.Println("解析到GenerateContentResponse失败", err)
		return nil, nil, err
	}
	newResp, err := c.ConvertChatResponse(&respBody)
	if err != nil {
		fmt.Println("转换到api.ChatCompletionNoStreamResponse失败", err)
		return nil, nil, err
	}
	newRespBytes, err := sonic.Marshal(newResp)
	if err != nil {
		return nil, nil, err
	}
	return io.NopCloser(bytes.NewReader(newRespBytes)), header, nil
}

func (c *Client) convertStreamChatCompletions(resp io.ReadCloser, header http.Header) (io.ReadCloser, http.Header, error) {
	pr, pw := io.Pipe()
	go func(w *io.PipeWriter) {
		defer w.Close()
		defer resp.Close()
		reader := bufio.NewReader(resp)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("读取数据失败", err)
				return
			}
			// 过滤掉可能的心跳包
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			// 解析 JSON 数据
			if string(line[:5]) == "data:" {
				var respBody GenerateContentResponse
				err = sonic.Unmarshal(line[5:], &respBody)
				if err != nil {
					fmt.Println("解析到GenerateContentResponse失败", err)
					return
				}
				newResp, err := c.ConvertStreamChatResponse(&respBody)
				if err != nil {
					fmt.Println("转换到api.ChatCompletionNoStreamResponse失败", err)
					return
				}
				newRespBytes, err := sonic.Marshal(newResp)
				if err != nil {
					fmt.Printf("data:%#v", line)
					fmt.Println("sonic.Marshal失败", err)
					return
				}
				data := fmt.Sprintf("data: %s\n\n", newRespBytes)
				_, err = w.Write([]byte(data))
				if err != nil {
					fmt.Println("w.Write失败", err)
					return
				}
			} else {
				fmt.Println("读取到的数据", string(line))
				_, err = w.Write(line)
				if err != nil {
					return
				}
			}
		}
	}(pw)
	return pr, header, nil
}

func (c *Client) ChatCompletions(ctx context.Context, req *api.ChatRequest) (io.ReadCloser, http.Header, error) {
	url := ""
	if req.Stream {
		baseURL := fmt.Sprintf("%s/v1beta/models/%s:streamGenerateContent", c.EndPoint, req.Model)
		params := stdurl.Values{}
		params.Set("key", c.APIKey)
		params.Set("alt", "sse")
		url = fmt.Sprintf("%s?%s", baseURL, params.Encode())
	} else {
		baseURL := fmt.Sprintf("%s/v1beta/models/%s:generateContent", c.EndPoint, req.Model)
		params := stdurl.Values{}
		params.Set("key", c.APIKey)
		url = fmt.Sprintf("%s?%s", baseURL, params.Encode())
	}
	newReqBody, err := c.ConvertChatRequest(req)
	if err != nil {
		return nil, nil, err
	}
	bodyBytes, err := sonic.Marshal(newReqBody)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(string(bodyBytes))
	resp, header, err := c.DoJsonRequest(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return resp, header, err
	}
	if req.Stream {
		return c.convertStreamChatCompletions(resp, header)
	} else {
		return c.convertNoStreamChatCompletions(resp, header)
	}
}

// ChatCompletionsByBytes 应该是没什么用的
func (c *Client) ChatCompletionsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	var r api.ChatRequest
	err := sonic.Unmarshal(req, &r)
	if err != nil {
		return nil, nil, err
	}
	return c.ChatCompletions(ctx, &r)
}

func content2Part(content []byte) ([]Part, error) {
	parts := make([]Part, 0)
	part := new(Part)
	var stringContent string
	err := json.Unmarshal(content, &stringContent)
	if err == nil {
		part.Text = stringContent
		parts = append(parts, *part)
		return parts, nil
	}
	var mediaContents []api.MediaContent
	err = json.Unmarshal(content, &mediaContents)
	if err == nil {
		for _, mediaContent := range mediaContents {
			switch mediaContent.Type {
			case "text":
				part.Text = mediaContent.Text
				parts = append(parts, *part)
				continue
			case "image_url":
				if mediaContent.ImageUrl.Url == "" {
					return nil, errors.New("image_url is empty")
				}
				part.FileData = &FileData{
					FileURL: mediaContent.ImageUrl.Url,
				}
				parts = append(parts, *part)
				continue
			case "input_audio":
				if mediaContent.InputAudio.Data == "" {
					return nil, errors.New("audio_data is empty")
				}
				if mediaContent.InputAudio.Format == "" {
					return nil, errors.New("audio_format is empty")
				}
				part.InlineData = &Blob{
					MimeType: "audio/" + mediaContent.InputAudio.Format,
					Data:     mediaContent.InputAudio.Data,
				}
				parts = append(parts, *part)
				continue
			default:
				return nil, errors.New("unknown content type")
			}
		}
	}
	return parts, err
}
