package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/jiu-u/oai-adapter/api"
	"github.com/jiu-u/oai-adapter/clients/legacy/openai"
	"github.com/jiu-u/oai-adapter/pkg/rr"
	"io"
	"log"
	"net/http"
	stdurl "net/url"
	"time"
)

type Client struct {
	*openai.Client
}

func NewClient(endpoint, apiKey string, proxy *stdurl.URL) *Client {
	return &Client{
		Client: openai.NewClient(endpoint, apiKey, proxy),
	}
}

func (c *Client) Completions(ctx context.Context, req *api.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	newReq := &OllamaGenerateRequest{
		Model:  req.Model,
		Prompt: req.Prompt,
		Suffix: req.Suffix,
		Images: req.Images,
		Format: req.Format,
		Options: Options{
			Temperature:      req.Temperature,
			Seed:             req.Seed,
			Echo:             req.Echo,
			FrequencyPenalty: req.FrequencyPenalty,
			Logprobs:         req.Logprobs,
			MaxTokens:        req.MaxTokens,
			BestOf:           req.BestOf,
			LogitBias:        req.LogitBias,
			N:                req.N,
			PresencePenalty:  req.PresencePenalty,
			Stop:             req.Stop,
			TopP:             req.TopP,
			User:             req.User,
		},
		System:    req.System,
		Template:  req.Template,
		Raw:       req.Raw,
		KeepAlive: req.KeepAlive,
		Stream:    req.Stream,
	}
	bodyBytes, err := sonic.Marshal(newReq)
	if err != nil {
		return nil, nil, err
	}
	url := c.EndPoint + "/api/generate"
	respReader, respHeader, err := c.DoJsonRequest(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return respReader, respHeader, err
	}
	if req.Stream {
		body, err := c.HandleCompletionsStream(respReader)
		header := new(http.Header)
		header.Set("Content-Type", "text/event-stream")
		header.Set("Cache-Control", "no-cache")
		header.Set("Connection", "keep-alive")
		header.Set("Transfer-Encoding", "chunked")
		return body, *header, err
	} else {
		body, err := c.HandleCompletionsNoStream(respReader)
		return body, respHeader, err
	}

}

func (c *Client) HandleCompletionsStream(resp io.ReadCloser) (io.ReadCloser, error) {
	r, w := io.Pipe()
	go func() {
		defer func() {
			err := resp.Close()
			if err != nil {
				fmt.Println("Error closing resp:", err)
			}
			err = w.Close()
			if err != nil {
				fmt.Println("Error closing w:", err)
			}
		}()
		scanner := bufio.NewScanner(resp)
		for scanner.Scan() {
			line := scanner.Text()

			// 解析 JSON 到结构体
			var lineData OllamaGenerateResponse
			if err := json.Unmarshal([]byte(line), &lineData); err != nil {
				log.Println("Error decoding JSON:", err)
				continue
			} else {
				openaiCompletions := c.GenResp2OpenAI(&lineData)
				respBytes, err := sonic.Marshal(openaiCompletions)
				if err != nil {
					log.Println("Error encoding respBytes:", err)
					continue
				}
				fmt.Fprintf(w, "data: %s\n\n", respBytes)
			}
		}
		fmt.Fprintf(w, "data: [DONE]\n\n")
	}()
	return r, nil
}

func (c *Client) GenResp2OpenAI(rawResp *OllamaGenerateResponse) *api.CompletionsResp {
	openaiCompletions := &api.CompletionsResp{
		ID:      "cmpl-" + rr.GenString(26),
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   rawResp.Model,
		Choices: []api.CompletionsChoice{
			{
				Text:         rawResp.Response,
				Index:        0,
				Logprobs:     nil,
				FinishReason: "stop",
			},
		},
		Usage: api.Usage{
			CompletionTokens:        rawResp.EvalCount,
			PromptTokens:            rawResp.PromptEvalCount,
			TotalTokens:             rawResp.EvalCount + rawResp.PromptEvalCount,
			CompletionTokensDetails: nil,
			PromptTokensDetails:     nil,
		},
		SystemFingerprint: "fp-" + rr.GenString(10),
	}
	return openaiCompletions
}

func (c *Client) HandleCompletionsNoStream(resp io.ReadCloser) (io.ReadCloser, error) {
	defer resp.Close()

	var ollamaResp OllamaGenerateResponse
	var err error
	//data,err := io.ReadAll(resp)
	//if err!= nil {
	//	return resp, err
	//}
	//err = sonic.Unmarshal(data, &ollamaResp)
	err = json.NewDecoder(resp).Decode(&ollamaResp)
	if err != nil {
		return nil, err
	}
	openaiCompletions := c.GenResp2OpenAI(&ollamaResp)
	respBytes, err := sonic.Marshal(openaiCompletions)
	if err != nil {
		return nil, err
	}
	resp = io.NopCloser(bytes.NewReader(respBytes))
	return resp, nil
}
