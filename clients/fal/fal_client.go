package fal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	v1 "github.com/jiu-u/oai-adapter/api/v1"
	"github.com/jiu-u/oai-adapter/clients/base"
	"github.com/jiu-u/oai-adapter/common"
	"github.com/jiu-u/oai-adapter/constant"
	"github.com/jiu-u/oai-adapter/pkg/task"
	"github.com/jiu-u/oai-adapter/tools"
	"github.com/oklog/ulid/v2"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Client struct {
	*base.Client
}

func NewClient(endPoint, apiKey string) *Client {
	if endPoint == "" {
		endPoint = constant.FalDefaultURL
	}
	endPoint = strings.TrimSpace(endPoint)
	endPoint = strings.TrimRight(endPoint, "/")
	return &Client{
		Client: base.NewClient(endPoint, apiKey),
	}
}

func (c *Client) SetHeader(header http.Header) {
	header.Del("Authorization")
	header.Set("Authorization", "Key "+c.APIKey)
}

func (c *Client) RelayRequest(ctx context.Context, method, targetPath string, body io.Reader, header http.Header) (io.ReadCloser, http.Header, error) {
	c.SetHeader(header)
	targetUrl := c.HomeUrl + targetPath
	return base.Relay(ctx, method, targetUrl, body, header, c.Client.Client)
}

func (c *Client) GenerateHeaderByContentType(contentType string) http.Header {
	headers := http.Header{}
	if len(contentType) > 0 {
		headers.Set("Content-Type", contentType)
	}
	headers.Set("Authorization", "Key "+c.APIKey)
	return headers
}

func (c *Client) CreateResponses(ctx context.Context, req *v1.ResponsesRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) ConvertChatCompletions(req *v1.ChatCompletionRequest) (io.ReadCloser, error) {
	var err error
	var req2 LLMRequest
	req2.Reasoning = true
	req2.Model = req.Model

	const maxChars = 4850

	var systemPromptText string
	var promptText string

	limitText := "请严格按照要求：你只能完成关于智能体补充或对user的回答，不可以扮演user，你的回答中不应该带有user。请严格按照要求！！！"
	systemPromptText += fmt.Sprintf("%s\n\n", limitText)

	for _, msg := range req.Messages {

		var content string
		if msg.IsStringContent() {
			content = msg.StringContent()
		} else {
			mediaContents, err := msg.ParseContent()
			if err != nil {
				return nil, fmt.Errorf("parse content error: %w", err)
			}
			for _, mediaContent := range mediaContents {
				// 对于其他类型内容，我们暂时忽略图像、音频等非文本类型
				if mediaContent.Type == v1.ContentTypeText {
					content += mediaContent.Text
				}
			}
		}
		formattedMsg := fmt.Sprintf("%s:%s\n\n", msg.Role, content)

		if utf8.RuneCountInString(formattedMsg)+utf8.RuneCountInString(systemPromptText) <= maxChars {
			systemPromptText += formattedMsg
			continue
		}
		if utf8.RuneCountInString(formattedMsg)+utf8.RuneCountInString(promptText) <= maxChars {
			promptText += formattedMsg
			continue
		}

	}

	req2.SystemPrompt = systemPromptText
	req2.Prompt = promptText

	// 确保Prompt不为空
	if req2.SystemPrompt == "" && req2.Prompt == "" {
		return nil, fmt.Errorf("prompt cannot be empty")
	}

	// 如果需要可以在这里删除末尾多余的换行符
	req2.Prompt = strings.TrimSuffix(req2.Prompt, "\n")
	req2.SystemPrompt = strings.TrimSuffix(req2.SystemPrompt, "\n")

	bytesData, err := sonic.Marshal(req2)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	rc := io.NopCloser(bytes.NewBuffer(bytesData))
	return rc, nil
}

func (c *Client) ConvertStreamChatCompletionsResponse(resp io.ReadCloser, model string) (io.ReadCloser, error) {
	r, w := io.Pipe()
	go func(w *io.PipeWriter) {
		defer func() {
			_ = resp.Close()
			_ = w.Close()
		}()

		scanner := bufio.NewScanner(resp)
		idx := 0
		for scanner.Scan() {
			line := scanner.Text()
			chunk := line

			// 过滤掉可能的心跳包
			if strings.TrimSpace(chunk) == "" || strings.HasPrefix(chunk, "event: event") {
				continue
			}

			// 解析 JSON 数据
			if strings.HasPrefix(line, "data:") {
				dataContent := strings.TrimSpace(line[5:])

				var chatResp v1.ChatCompletionStreamResponse
				var llmResp LLMResponse

				err := json.Unmarshal([]byte(dataContent), &llmResp)

				if err != nil {
					fmt.Println("Error decoding JSON:", err)
					continue
				} else {
					msgText := llmResp.Output[idx:]
					idx = len(llmResp.Output)
					chatResp = v1.ChatCompletionStreamResponse{
						Choices: []v1.ChoiceWithDelta{
							{
								Delta: v1.Delta{
									Content: msgText,
								},
								FinishReason: "",
								Index:        0,
								Logprobs:     nil,
							},
						},
						Created:           time.Now().Unix(),
						ID:                "chatcmpl-" + ulid.Make().String(),
						Model:             model,
						Object:            "chat.completion.chunk",
						SystemFingerprint: "fp_" + ulid.Make().String(),
						Usage:             nil,
					}
					respBytes, err := sonic.Marshal(chatResp)
					if err != nil {
						fmt.Printf("data:%#v", line)
						fmt.Println("sonic.Marshal失败", err)
						return
					}
					_, err = fmt.Fprintf(w, "data: %s\n\n", respBytes)

					if err != nil {
						fmt.Println("Error writing to pipe:", err)
						return
					}
				}
			}
		}

		_, _ = fmt.Fprintf(w, "data: [DONE]\n\n")

		// 处理扫描器可能的错误
		if err := scanner.Err(); err != nil {
			fmt.Println("Scanner error:", err)
		}
		time.Sleep(time.Second * 1)

	}(w)

	return r, nil
}

func (c *Client) CreateChatCompletions(ctx context.Context, req *v1.ChatCompletionRequest) (io.ReadCloser, http.Header, error) {
	var err error
	reqBody, err := c.ConvertChatCompletions(req)
	if err != nil {
		return nil, nil, err
	}
	targetUrl := HomeUrl + "/fal-ai/any-llm/stream"
	if req.Stream {
		body, _, err := base.Relay(ctx, http.MethodPost, targetUrl, reqBody, c.GenerateHeaderByContentType("application/json"), c.Client.Client)
		if err != nil {
			return nil, nil, err
		}
		rc, err := c.ConvertStreamChatCompletionsResponse(body, req.Model)
		if err != nil {
			return nil, nil, err
		}
		header := http.Header{}
		header.Set("Content-Type", "text/event-stream")
		header.Set("Cache-Control", "no-cache")
		header.Set("Connection", "keep-alive")
		header.Set("Transfer-Encoding", "chunked")
		return rc, header, nil
	} else {
		body, _, err := base.Relay(ctx, http.MethodPost, targetUrl, reqBody, c.GenerateHeaderByContentType("application/json"), c.Client.Client)
		//body, _, err := c.RelayRequest(ctx, http.MethodPost, "/fal-ai/any-llm/stream", reqBody, c.GenerateHeaderByContentType("application/json"))
		if err != nil {
			return nil, nil, err
		}
		var data LLMResponse
		scanner := bufio.NewScanner(body)
		for scanner.Scan() {
			line := scanner.Bytes()
			chunk := string(line)
			if strings.TrimSpace(chunk) == "" || strings.HasPrefix(chunk, "event: event") {
				continue
			}
			if !strings.HasPrefix(chunk, "data: ") {
				continue
			}
			chunk = strings.TrimPrefix(chunk, "data: ")
			chunk = strings.TrimSpace(chunk)
			err := json.Unmarshal([]byte(chunk), &data)
			if err != nil {
				return nil, nil, err
			}
			if data.Partial == false {
				break
			}
		}

		resp := v1.ChatCompletionResponse{
			ID:      "chatcmpl-" + ulid.Make().String(),
			Model:   req.Model,
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Choices: []v1.Choice{
				{
					Index:        0,
					FinishReason: "stop",
					Message: v1.CompletionMessage{
						Role:    "assistant",
						Content: data.Output,
					},
				},
			},
		}
		resBytes, err := sonic.Marshal(resp)
		if err != nil {
			return nil, nil, err
		}
		rc := io.NopCloser(bytes.NewBuffer(resBytes))
		header := http.Header{}
		header.Set("Content-Type", "application/json")
		header.Set("Cache-Control", "no-cache")
		return rc, header, nil
	}

}

func (c *Client) CreateCompletions(ctx context.Context, req *v1.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) Models(ctx context.Context) (*v1.ModelResponse, error) {
	resp := v1.ModelResponse{
		Object: "list",
		Data:   make([]v1.Model, 0, len(Models)),
	}
	for _, model := range Models {
		resp.Data = append(resp.Data, v1.Model{
			ID:      model,
			Object:  "model",
			Created: 0,
			OwnedBy: "fal",
		})
	}
	return &resp, nil
}

func (c *Client) CreateEmbeddings(ctx context.Context, req *v1.EmbeddingsRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateRerank(ctx context.Context, req *v1.RerankRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateImage(ctx context.Context, req *v1.ImageGenerateRequest) (io.ReadCloser, http.Header, error) {
	imageReq := ImageCreateRequest{
		Prompt:            req.Prompt,
		ImageSize:         nil,
		NegativePrompt:    req.NegativePrompt,
		Seed:              req.Seed,
		NumInferenceSteps: req.NumInferenceSteps,
		NumImages:         req.N,
	}
	if req.Size != "" {
		imageSize := strings.Split(req.Size, "x")
		if len(imageSize) == 2 {
			width, _ := strconv.ParseInt(imageSize[0], 10, 64)
			height, _ := strconv.ParseInt(imageSize[1], 10, 64)
			imageReq.ImageSize = &ImageSize{
				Width:  int(width),
				Height: int(height),
			}
		}
	}
	reqBytes, err := json.Marshal(imageReq)
	if err != nil {
		return nil, nil, err
	}

	reqBody := bytes.NewBuffer(reqBytes)
	// 发送请求
	createPath := "/fal-ai/" + strings.TrimSpace(req.Model)

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	c.SetHeader(header)

	createResp, _, err := c.RelayRequest(ctx, http.MethodPost, createPath, reqBody, header)
	if err != nil {
		return nil, nil, err
	}

	createRespBytes, _ := io.ReadAll(createResp)
	var queueResp QueueResponse
	err = sonic.Unmarshal(createRespBytes, &queueResp)
	if err != nil {
		return nil, nil, err
	}
	respBytes, err := c.DoPollTask(ctx, &queueResp)
	if err != nil {
		return nil, nil, err
	}
	var resp ImageCreateResponse
	err = sonic.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, nil, err
	}
	openaiImageResp := v1.ImageGenerateResponse{
		Created: time.Now().Unix(),
	}
	for _, image := range resp.Images {
		openaiImageResp.Data = append(openaiImageResp.Data, v1.ImageGenData{
			URL:           image.URL,
			B64JSON:       "",
			RevisedPrompt: resp.Prompt,
		})
	}

	oaiImageRespBytes, err := sonic.Marshal(openaiImageResp)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal error: %w", err)
	}
	respHeader := http.Header{}
	respHeader.Set("Content-Type", "application/json")
	return io.NopCloser(bytes.NewBuffer(oaiImageRespBytes)), respHeader, nil
}

// 计算下一个间隔时间，使用指数退避策略
func (c *Client) calculateNextInterval(currentInterval time.Duration, options *task.PollTaskOptions) time.Duration {
	nextInterval := time.Duration(float64(currentInterval) * options.BackoffFactor)
	if nextInterval > options.MaxInterval {
		return options.MaxInterval
	}
	return nextInterval
}

// DoPollTask 轮询任务状态并获取结果
func (c *Client) DoPollTask(ctx context.Context, queueResp *QueueResponse, opts ...*task.PollTaskOptions) ([]byte, error) {
	// 使用提供的选项或默认值
	var options *task.PollTaskOptions
	if len(opts) > 0 {
		options = opts[0]
	} else {
		options = &common.DefaultPollTaskOptions
	}

	// 创建带超时的上下文
	var cancelFunc context.CancelFunc
	if options.Timeout > 0 {
		ctx, cancelFunc = context.WithTimeout(ctx, options.Timeout)
		defer cancelFunc()
	}

	interval := options.InitialInterval
	var lastErr error
	time.Sleep(interval)

	for attempt := 0; attempt < options.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			if lastErr != nil {
				return nil, fmt.Errorf("polling canceled: %w (last error: %v)", ctx.Err(), lastErr)
			}
			return nil, fmt.Errorf("polling canceled: %w", ctx.Err())
		default:
			// 继续执行
		}

		// 检查任务状态
		statusResp, _, err := base.Relay(ctx, http.MethodGet, queueResp.StatusUrl, nil,
			c.GenerateHeaderByContentType(""), c.Client.Client)

		if err != nil {
			lastErr = fmt.Errorf("status request failed: %w", err)
			//c.logPollError(queueResp, attempt, lastErr)
			time.Sleep(interval)
			interval = c.calculateNextInterval(interval, options)
			continue
		}

		// 解析状态响应
		statusRespBytes, err := io.ReadAll(statusResp)
		if err != nil {
			lastErr = fmt.Errorf("failed to read status response: %w", err)
			//c.logPollError(queueResp, attempt, lastErr)
			time.Sleep(interval)
			interval = c.calculateNextInterval(interval, options)
			continue
		}

		var queueStatus QueueResponse
		if err = sonic.Unmarshal(statusRespBytes, &queueStatus); err != nil {
			lastErr = fmt.Errorf("failed to parse status response: %w", err)
			//c.logPollError(queueResp, attempt, lastErr)
			time.Sleep(interval)
			interval = c.calculateNextInterval(interval, options)
			continue
		}

		// 根据任务状态处理
		switch queueStatus.Status {
		case "COMPLETED":
			header := http.Header{}
			header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
			header.Add("Content-Type", "application/json")
			header.Add("Accept", "*/*")

			c.SetHeader(header)
			// 任务完成，获取结果
			outputResp, _, err := base.Relay(ctx, http.MethodGet, queueResp.ResponseUrl, nil,
				header, c.Client.Client)

			if err != nil {
				if errors.Is(err, common.InternalError) {

					lastErr = fmt.Errorf("failed to fetch completed task result: %v", err)
				} else {
					lastErr = fmt.Errorf("failed to fetch completed task result: %v", err)
				}
				time.Sleep(interval)
				interval = c.calculateNextInterval(interval, options)
				continue
			}

			outputRespBytes, err := io.ReadAll(outputResp)

			if err != nil {
				lastErr = fmt.Errorf("failed to read task result: %v", err)
				time.Sleep(interval)
				interval = c.calculateNextInterval(interval, options)
				continue
			}

			return outputRespBytes, nil

		case "FAILED":
			// 任务失败
			var errMsg string
			if queueStatus.Logs != nil {
				errMsg = fmt.Sprintf("task failed with logs: %s", string(queueStatus.Logs))
			} else {
				errMsg = "task failed without error details"
			}
			return nil, fmt.Errorf(errMsg)

		case "CANCELED":
			return nil, fmt.Errorf("task was canceled")

		default:
			// 任务仍在处理中，继续轮询
			//c.logTaskProgress(queueResp, attempt, queueStatus.Status)
			time.Sleep(interval)
			interval = c.calculateNextInterval(interval, options)
		}
	}

	// 达到最大尝试次数
	if lastErr != nil {
		return nil, fmt.Errorf("polling exceeded maximum attempts: %w", lastErr)
	}
	return nil, fmt.Errorf("polling exceeded maximum attempts (%d) without completion", options.MaxAttempts)
}

func (c *Client) CreateImageEdit(ctx context.Context, req *v1.ImageEditRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateImageVariation(ctx context.Context, req *v1.ImageVariationRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateSpeech(ctx context.Context, req *v1.AudioSpeechRequest) (io.ReadCloser, http.Header, error) {
	createRequest := Text2AudioRequest{
		Prompt: req.Input,
		Voice:  req.Voice,
		Speed:  req.Speed,
	}

	reqBytes, err := json.Marshal(createRequest)
	if err != nil {
		return nil, nil, err
	}

	reqBody := bytes.NewBuffer(reqBytes)
	// 发送请求
	createPath := "/fal-ai/" + strings.TrimSpace(req.Model)

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	c.SetHeader(header)

	createResp, _, err := c.RelayRequest(ctx, http.MethodPost, createPath, reqBody, header)
	if err != nil {
		return nil, nil, err
	}

	createRespBytes, _ := io.ReadAll(createResp)
	var queueResp QueueResponse
	err = sonic.Unmarshal(createRespBytes, &queueResp)
	if err != nil {
		return nil, nil, err
	}
	respBytes, err := c.DoPollTask(ctx, &queueResp)
	if err != nil {
		return nil, nil, err
	}
	var resp Text2AudioResponse
	err = sonic.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, nil, err
	}
	rc, respHeader, err := tools.GetReadCloserFromURL(resp.Audio.URL)
	if err != nil {
		return nil, nil, err
	}
	return rc, respHeader, nil
}

func (c *Client) CreateTranslation(ctx context.Context, req *v1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateTranscription(ctx context.Context, req *v1.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateVideoSubmit(ctx context.Context, req *v1.VideoRequest) (*v1.VideoResponse, error) {
	createRequest := TextToVideoRequest{
		Prompt:         req.Prompt,
		Resolution:     req.ImageSize,
		NegativePrompt: req.NegativePrompt,
		Style:          req.Style,
		Seed:           req.Seed,
		ImageUrl:       req.Image,
	}

	reqBytes, err := json.Marshal(createRequest)
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewBuffer(reqBytes)
	// 发送请求
	createPath := "/fal-ai/" + strings.TrimSpace(req.Model)

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	c.SetHeader(header)

	createResp, _, err := c.RelayRequest(ctx, http.MethodPost, createPath, reqBody, header)
	if err != nil {
		return nil, err
	}

	createRespBytes, _ := io.ReadAll(createResp)
	var queueResp QueueResponse
	err = sonic.Unmarshal(createRespBytes, &queueResp)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(context.Background(), "model", req.Model)
	// task
	poller := base.NewPoller(c.GetVideoStatus)
	taskID := c.TaskMgr.CreatePollingTask(ctx, queueResp.RequestId, poller, nil)

	resp := &v1.VideoResponse{
		RequestId: taskID,
	}
	resp.RequestId = taskID
	return resp, nil

}

func (c *Client) GetVideoStatus(ctx context.Context, externalID string) (bool, any, error) {
	var videoResp v1.VideoStatusResponse
	videoResp.Status = "InQueue"
	model := ctx.Value("model").(string)
	startPath := "/fal-ai/" + strings.TrimSpace(model) + "/requests"

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	c.SetHeader(header)

	statusUrl := c.HomeUrl + startPath + "/" + externalID + "/status"

	statusResp, _, err := base.Relay(ctx, http.MethodGet, statusUrl, nil,
		c.GenerateHeaderByContentType(""), c.Client.Client)

	if err != nil {
		return false, &videoResp, err
	}

	// 解析状态响应
	statusRespBytes, err := io.ReadAll(statusResp)
	if err != nil {
		return false, &videoResp, err
	}

	var queueStatus QueueResponse
	if err = sonic.Unmarshal(statusRespBytes, &queueStatus); err != nil {
		return false, &videoResp, err
	}
	// 根据任务状态处理
	responseUrl := c.HomeUrl + startPath + "/" + externalID
	switch queueStatus.Status {
	case "COMPLETED":

		header := http.Header{}
		header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
		header.Add("Content-Type", "application/json")
		header.Add("Accept", "*/*")

		c.SetHeader(header)
		// 任务完成，获取结果
		outputResp, _, err := base.Relay(ctx, http.MethodGet, responseUrl, nil,
			header, c.Client.Client)

		if err != nil {
			return false, &videoResp, err
		}

		outputRespBytes, err := io.ReadAll(outputResp)

		if err != nil {
			return false, &videoResp, err
		}

		var v TextToVideoResponse
		err = json.Unmarshal(outputRespBytes, &v)
		if err != nil {
			return false, &videoResp, err
		}
		videoResp.Status = "Succeed"
		videoResp.Results = []v1.VideoResult{
			{
				Videos: []v1.VideoItem{
					{
						Url: v.Video.Url,
					},
				},
			},
		}
		return true, &videoResp, nil

	case "FAILED":
		videoResp.Status = "Failed"
		videoResp.Reason = "task failed"
		return true, &videoResp, nil

	case "CANCELED":
		videoResp.Status = "Failed"
		videoResp.Reason = "task was canceled"
		return true, &videoResp, nil

	default:
		videoResp.Status = "InProgress"
		videoResp.Reason = "task is still processing"
		return false, &videoResp, nil
	}

}
