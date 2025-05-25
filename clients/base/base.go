package base

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	v1 "github.com/jiu-u/oai-adapter/api/v1"
	"github.com/jiu-u/oai-adapter/common"
	"github.com/jiu-u/oai-adapter/pkg/task"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	EndPoint string // 使用v1后缀 eg: https://api.openai.com/v1
	APIKey   string
	Client   *http.Client
	TaskMgr  task.TaskManager
	HomeUrl  string
}

func NewClient(EndPoint, apiKey string) *Client {
	EndPoint = strings.TrimSpace(EndPoint)
	EndPoint = strings.TrimRight(EndPoint, "/")
	parsedURL, _ := url.Parse(EndPoint)
	baseUrl := parsedURL.Scheme + "://" + parsedURL.Host
	return &Client{
		EndPoint: EndPoint,
		APIKey:   apiKey,
		Client:   common.GetDefaultClient(),
		TaskMgr:  common.GetDefaultTaskManager(),
		HomeUrl:  baseUrl,
	}
}

func (c *Client) SetClient(client *http.Client) {
	c.Client = client
}

func (c *Client) RelayRequest(ctx context.Context, method, targetPath string, body io.Reader, header http.Header) (io.ReadCloser, http.Header, error) {
	c.SetHeader(header)
	targetUrl := c.HomeUrl + targetPath
	return Relay(ctx, method, targetUrl, body, header, c.Client)
}

func (c *Client) SetHeader(header http.Header) {
	header.Del("Authorization")
	header.Set("Authorization", "Bearer "+c.APIKey)
}

func (c *Client) GenerateHeaderByContentType(contentType string) http.Header {
	headers := http.Header{}
	headers.Set("Content-Type", contentType)
	headers.Set("Authorization", "Bearer "+c.APIKey)
	return headers
}

//func (c *Client) ConvertChatCompletions(req *v1.ChatCompletionRequest)(io.ReadCloser, http.Header, error)
// completions
// embeddings
// rerank
// imageGen
// imageEdit
// imageVariation
// speech
// translation
// transcription
// videoSubmit
// videoStatus

func (c *Client) SamePostJob(ctx context.Context, targetUrl string, req any, contentType string) (io.ReadCloser, http.Header, error) {
	var err error
	reqBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal error: %w", err)
	}
	body := bytes.NewBuffer(reqBytes)
	header := c.GenerateHeaderByContentType(contentType)
	return Relay(ctx, http.MethodPost, targetUrl, body, header, c.Client)
}

func (c *Client) CreateResponses(ctx context.Context, req *v1.ResponsesRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/responses"
	return c.SamePostJob(ctx, targetUrl, req, "application/json")

}

func (c *Client) CreateChatCompletions(ctx context.Context, req *v1.ChatCompletionRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/chat/completions"
	return c.SamePostJob(ctx, targetUrl, req, "application/json")
}

func (c *Client) CreateCompletions(ctx context.Context, req *v1.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/completions"
	return c.SamePostJob(ctx, targetUrl, req, "application/json")
}

func (c *Client) Models(ctx context.Context) (*v1.ModelResponse, error) {
	var err error
	targetUrl := c.EndPoint + "/models"
	data, _, err := Relay(ctx, http.MethodGet, targetUrl, nil, nil, c.Client)
	if err != nil {
		return nil, err
	}
	dataBytes, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("read all error: %w", err)
	}
	var resp v1.ModelResponse
	err = sonic.Unmarshal(dataBytes, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return &resp, nil
}

func (c *Client) CreateEmbeddings(ctx context.Context, req *v1.EmbeddingsRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/embeddings"
	return c.SamePostJob(ctx, targetUrl, req, "application/json")
}

func (c *Client) CreateSpeech(ctx context.Context, req *v1.AudioSpeechRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/audio/speech"
	return c.SamePostJob(ctx, targetUrl, req, "application/json")
}

func (c *Client) CreateTranslation(ctx context.Context, req *v1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/audio/translations"

	// 创建一个字节缓冲区来存储请求体
	var buf bytes.Buffer
	// 创建 multipart 写入器
	writer := multipart.NewWriter(&buf)

	// 添加文件字段
	part, err := writer.CreateFormFile("file", req.File.Filename)
	if err != nil {
		return nil, nil, fmt.Errorf("create form file error: %w", err)
	}

	file, err := req.File.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("open file error: %w", err)
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 将文件内容复制到文件字段
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, fmt.Errorf("copy file content error: %w", err)
	}

	// 添加其他字段
	if err := writer.WriteField("model", req.Model); err != nil {
		return nil, nil, fmt.Errorf("write model field error: %w", err)
	}

	if req.Prompt != "" {
		if err := writer.WriteField("prompt", req.Prompt); err != nil {
			return nil, nil, fmt.Errorf("write prompt field error: %w", err)
		}
	}

	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", req.ResponseFormat); err != nil {
			return nil, nil, fmt.Errorf("write response_format field error: %w", err)
		}
	}

	if req.Temperature != 0 {
		if err := writer.WriteField("temperature", strconv.FormatFloat(req.Temperature, 'f', -1, 64)); err != nil {
			return nil, nil, fmt.Errorf("write temperature field error: %w", err)
		}
	}

	// 关闭 multipart 写入器
	err = writer.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("close multipart writer error: %w", err)
	}

	// 生成请求头
	header := http.Header{}
	header.Set("Content-Type", writer.FormDataContentType())
	header.Set("Authorization", "Bearer "+c.APIKey)

	// 返回请求
	return Relay(ctx, http.MethodPost, targetUrl, &buf, header, c.Client)
}

func (c *Client) CreateTranscription(ctx context.Context, req *v1.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/audio/transcriptions"

	// 创建一个字节缓冲区来存储请求体
	var buf bytes.Buffer
	// 创建 multipart 写入器
	writer := multipart.NewWriter(&buf)

	// 添加文件字段
	part, err := writer.CreateFormFile("file", req.File.Filename)
	if err != nil {
		return nil, nil, fmt.Errorf("create form file error: %w", err)
	}

	file, err := req.File.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("open file error: %w", err)
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 将文件内容复制到文件字段
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, fmt.Errorf("copy file content error: %w", err)
	}

	// 添加必填字段
	if err := writer.WriteField("model", req.Model); err != nil {
		return nil, nil, fmt.Errorf("write model field error: %w", err)
	}

	// 添加可选字段
	if req.ChunkingStrategy != "" {
		if err := writer.WriteField("chunking_strategy", req.ChunkingStrategy); err != nil {
			return nil, nil, fmt.Errorf("write chunking_strategy field error: %w", err)
		}
	}

	// 处理包含数组
	if len(req.Include) > 0 {
		for _, item := range req.Include {
			if err := writer.WriteField("include", item); err != nil {
				return nil, nil, fmt.Errorf("write include field error: %w", err)
			}
		}
	}

	if req.Language != "" {
		if err := writer.WriteField("language", req.Language); err != nil {
			return nil, nil, fmt.Errorf("write language field error: %w", err)
		}
	}

	if req.Prompt != "" {
		if err := writer.WriteField("prompt", req.Prompt); err != nil {
			return nil, nil, fmt.Errorf("write prompt field error: %w", err)
		}
	}

	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", req.ResponseFormat); err != nil {
			return nil, nil, fmt.Errorf("write response_format field error: %w", err)
		}
	}

	if req.Temperature != 0 {
		if err := writer.WriteField("temperature", strconv.FormatFloat(req.Temperature, 'f', -1, 64)); err != nil {
			return nil, nil, fmt.Errorf("write temperature field error: %w", err)
		}
	}

	// 处理时间戳粒度数组
	if len(req.TimestampGranularities) > 0 {
		for _, item := range req.TimestampGranularities {
			if err := writer.WriteField("timestamp_granularities", item); err != nil {
				return nil, nil, fmt.Errorf("write timestamp_granularities field error: %w", err)
			}
		}
	}

	// 关闭 multipart 写入器
	err = writer.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("close multipart writer error: %w", err)
	}

	// 生成请求头
	header := http.Header{}
	header.Set("Content-Type", writer.FormDataContentType())
	header.Set("Authorization", "Bearer "+c.APIKey)

	// 返回请求
	return Relay(ctx, http.MethodPost, targetUrl, &buf, header, c.Client)
}

func (c *Client) CreateImage(ctx context.Context, req *v1.ImageGenerateRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/images/generations"
	return c.SamePostJob(ctx, targetUrl, req, "application/json")
}

func (c *Client) CreateImageEdit(ctx context.Context, req *v1.ImageEditRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/images/edits"

	// 创建一个字节缓冲区来存储请求体
	var buf bytes.Buffer
	// 创建 multipart 写入器
	writer := multipart.NewWriter(&buf)

	// 添加图片字段，支持多个图片
	for i, image := range req.Image {
		part, err := writer.CreateFormFile(fmt.Sprintf("image[%d]", i), image.Filename)
		if err != nil {
			return nil, nil, fmt.Errorf("create form file for image error: %w", err)
		}

		file, err := image.Open()
		if err != nil {
			return nil, nil, fmt.Errorf("open image file error: %w", err)
		}
		defer file.Close()

		// 将文件内容复制到文件字段
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, nil, fmt.Errorf("copy image file content error: %w", err)
		}
	}

	// 添加遮罩文件（如果有）
	if req.Mask != nil {
		part, err := writer.CreateFormFile("mask", req.Mask.Filename)
		if err != nil {
			return nil, nil, fmt.Errorf("create form file for mask error: %w", err)
		}

		file, err := req.Mask.Open()
		if err != nil {
			return nil, nil, fmt.Errorf("open mask file error: %w", err)
		}
		defer file.Close()

		// 将文件内容复制到文件字段
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, nil, fmt.Errorf("copy mask file content error: %w", err)
		}
	}

	// 添加必填字段
	if err := writer.WriteField("prompt", req.Prompt); err != nil {
		return nil, nil, fmt.Errorf("write prompt field error: %w", err)
	}

	if err := writer.WriteField("model", req.Model); err != nil {
		return nil, nil, fmt.Errorf("write model field error: %w", err)
	}

	// 添加可选字段
	if req.Background != "" {
		if err := writer.WriteField("background", req.Background); err != nil {
			return nil, nil, fmt.Errorf("write background field error: %w", err)
		}
	}

	if req.N > 0 {
		if err := writer.WriteField("n", strconv.Itoa(req.N)); err != nil {
			return nil, nil, fmt.Errorf("write n field error: %w", err)
		}
	}

	if req.Quality != "" {
		if err := writer.WriteField("quality", req.Quality); err != nil {
			return nil, nil, fmt.Errorf("write quality field error: %w", err)
		}
	}

	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", req.ResponseFormat); err != nil {
			return nil, nil, fmt.Errorf("write response_format field error: %w", err)
		}
	}

	if req.Size != "" {
		if err := writer.WriteField("size", req.Size); err != nil {
			return nil, nil, fmt.Errorf("write size field error: %w", err)
		}
	}

	if req.User != "" {
		if err := writer.WriteField("user", req.User); err != nil {
			return nil, nil, fmt.Errorf("write user field error: %w", err)
		}
	}

	// 关闭 multipart 写入器
	err := writer.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("close multipart writer error: %w", err)
	}

	// 生成请求头
	header := http.Header{}
	header.Set("Content-Type", writer.FormDataContentType())
	header.Set("Authorization", "Bearer "+c.APIKey)

	// 返回请求
	return Relay(ctx, http.MethodPost, targetUrl, &buf, header, c.Client)
}

func (c *Client) CreateImageVariation(ctx context.Context, req *v1.ImageVariationRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/images/variations"

	// 创建一个字节缓冲区来存储请求体
	var buf bytes.Buffer
	// 创建 multipart 写入器
	writer := multipart.NewWriter(&buf)

	// 添加图片文件
	part, err := writer.CreateFormFile("image", req.Image.Filename)
	if err != nil {
		return nil, nil, fmt.Errorf("create form file for image error: %w", err)
	}

	file, err := req.Image.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("open image file error: %w", err)
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 将文件内容复制到文件字段
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, fmt.Errorf("copy image file content error: %w", err)
	}

	// 添加必填字段
	if err := writer.WriteField("model", req.Model); err != nil {
		return nil, nil, fmt.Errorf("write model field error: %w", err)
	}

	// 添加可选字段
	if req.N > 0 {
		if err := writer.WriteField("n", strconv.Itoa(req.N)); err != nil {
			return nil, nil, fmt.Errorf("write n field error: %w", err)
		}
	}

	if req.Size != "" {
		if err := writer.WriteField("size", req.Size); err != nil {
			return nil, nil, fmt.Errorf("write size field error: %w", err)
		}
	}

	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", req.ResponseFormat); err != nil {
			return nil, nil, fmt.Errorf("write response_format field error: %w", err)
		}
	}

	if req.User != "" {
		if err := writer.WriteField("user", req.User); err != nil {
			return nil, nil, fmt.Errorf("write user field error: %w", err)
		}
	}

	// 关闭 multipart 写入器
	err = writer.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("close multipart writer error: %w", err)
	}

	// 生成请求头
	header := http.Header{}
	header.Set("Content-Type", writer.FormDataContentType())
	header.Set("Authorization", "Bearer "+c.APIKey)

	// 返回请求
	return Relay(ctx, http.MethodPost, targetUrl, &buf, header, c.Client)
}

func (c *Client) CreateRerank(ctx context.Context, req *v1.RerankRequest) (io.ReadCloser, http.Header, error) {
	targetUrl := c.EndPoint + "/rerank"
	return c.SamePostJob(ctx, targetUrl, req, "application/json")
}

func (c *Client) CreateVideoSubmit(ctx context.Context, req *v1.VideoRequest) (*v1.VideoResponse, error) {
	var err error
	targetUrl := c.EndPoint + "/videos/submit"
	respBody, _, err := c.SamePostJob(ctx, targetUrl, req, "application/json")
	if err != nil {
		return nil, fmt.Errorf("relay error: %w", err)
	}
	respBodyBytes, err := io.ReadAll(respBody)
	if err != nil {
		return nil, fmt.Errorf("read all error: %w", err)
	}
	var resp v1.VideoResponse
	err = sonic.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	// task
	poller := NewPoller(c.GetVideoStatus)
	taskID := c.TaskMgr.CreatePollingTask(ctx, resp.RequestId, poller, nil)
	resp.RequestId = taskID
	return &resp, nil
}

func (c *Client) GetVideoStatus(ctx context.Context, externalID string) (bool, any, error) {
	var err error
	var resp v1.VideoStatusResponse
	resp.RawRequestId = externalID
	targetUrl := c.EndPoint + "/videos/status"
	header := c.GenerateHeaderByContentType("application/json")
	req := v1.VideoStatusRequest{
		RequestId: externalID,
	}
	bodyBytes, _ := sonic.Marshal(req)
	body := bytes.NewReader(bodyBytes)
	respBody, _, err := Relay(ctx, http.MethodPost, targetUrl, body, header, c.Client)
	if err != nil {
		resp.Status = "Failed"
		resp.Reason = "relay error"
		return false, &resp, fmt.Errorf("relay error: %w", err)
	}
	respBodyBytes, err := io.ReadAll(respBody)
	if err != nil {
		resp.Status = "Failed"
		resp.Reason = "read Data error"
		return false, &resp, fmt.Errorf("read all error: %w", err)
	}
	err = sonic.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		resp.Status = "Failed"
		resp.Reason = "unmarshal error"
		return false, &resp, fmt.Errorf("unmarshal error: %w", err)
	}
	if resp.Status == "Succeed" || resp.Status == "Failed" {
		return true, &resp, nil
	}
	return false, &resp, nil
}
