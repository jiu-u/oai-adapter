package volcengine

import (
	"bytes"
	"context"
	"github.com/bytedance/sonic"
	v1 "github.com/jiu-u/oai-adapter/api/v1"
	"github.com/jiu-u/oai-adapter/clients/base"
	"github.com/jiu-u/oai-adapter/constant"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	*base.Client
}

func NewClient(endPoint, apiKey string) *Client {
	if endPoint == "" {
		endPoint = constant.VolcEngineDefaultURL
	}
	endPoint = strings.TrimSpace(endPoint)
	endPoint = strings.TrimRight(endPoint, "/")
	endPoint = endPoint + "/api/v3"

	return &Client{
		Client: base.NewClient(endPoint, apiKey),
	}
}

func (c *Client) CreateResponses(ctx context.Context, req *v1.ResponsesRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateCompletions(ctx context.Context, req *v1.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateRerank(ctx context.Context, req *v1.RerankRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateSpeech(ctx context.Context, req *v1.AudioSpeechRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateTranslation(ctx context.Context, req *v1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateTranscription(ctx context.Context, req *v1.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateImageEdit(ctx context.Context, req *v1.ImageEditRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateImageVariation(ctx context.Context, req *v1.ImageVariationRequest) (io.ReadCloser, http.Header, error) {
	return base.NoImplementMethod(ctx, req)
}

func (c *Client) CreateVideoSubmit(ctx context.Context, req *v1.VideoRequest) (*v1.VideoResponse, error) {
	targetUrl := c.EndPoint + "/contents/generations/tasks"
	vReq := &CreateVideoTaskRequest{
		Model:   req.Model,
		Content: nil,
		Seed:    int64(req.Seed),
	}
	if req.Prompt != "" {
		vReq.Content = append(vReq.Content, CreateVideoTaskContent{
			Type:     "text",
			Text:     req.Prompt,
			ImageUrl: nil,
		})
	}
	if req.Image != "" {
		vReq.Content = append(vReq.Content, CreateVideoTaskContent{
			Type:     "image_url",
			Text:     "",
			ImageUrl: &CreateVideoTaskImage{Url: req.Image},
		})
	}
	reqBytes, err := sonic.Marshal(vReq)
	if err != nil {
		return nil, err
	}
	resp, _, err := base.Relay(ctx, http.MethodPost, targetUrl, bytes.NewBuffer(reqBytes), c.GenerateHeaderByContentType("application/json"), c.Client.Client)
	if err != nil {
		return nil, err
	}
	respBytes, err := io.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	var vResp VideoTaskStatusResponse
	err = sonic.Unmarshal(respBytes, &vResp)
	if err != nil {
		return nil, err
	}
	return &v1.VideoResponse{
		RequestId: vResp.Id,
	}, nil
}

func (c *Client) GetVideoStatus(ctx context.Context, externalID string) (bool, any, error) {
	var err error
	var resp v1.VideoStatusResponse
	resp.Status = "InProgress"
	resp.RawRequestId = externalID
	targetUrl := c.EndPoint + "/contents/generations/tasks/" + externalID
	vReq := &VideoTaskStatusRequest{
		Id: externalID,
	}
	reqBytes, err := sonic.Marshal(vReq)
	if err != nil {
		resp.Error = err
		return false, &resp, err
	}
	vResp, _, err := base.Relay(ctx, http.MethodGet, targetUrl, bytes.NewBuffer(reqBytes), c.GenerateHeaderByContentType("application/json"), c.Client.Client)
	if err != nil {
		resp.Error = err
		return false, &resp, err
	}
	vRespBytes, err := io.ReadAll(vResp)
	if err != nil {
		resp.Error = err
		return false, &resp, err
	}
	var vStatusResp VideoTaskStatusResponse
	err = sonic.Unmarshal(vRespBytes, &vStatusResp)
	if err != nil {
		resp.Error = err
		return false, &resp, err
	}
	if vStatusResp.Status == "succeeded" || vStatusResp.Status == "failed" {
		if vStatusResp.Status == "succeeded" {
			resp.Status = "Succeed"
		} else {
			resp.Status = "Failed"
		}
		result := v1.VideoResult{
			Videos: []v1.VideoItem{
				{
					Url: vStatusResp.Content.VideoUrl,
				},
			},
		}
		resp.Results = []v1.VideoResult{result}
		resp.Status = "Failed"
		resp.Error = vStatusResp.Error
		return true, &resp, nil
	}
	return false, &resp, err
}
