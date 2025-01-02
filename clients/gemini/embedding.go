package gemini

import (
	"bytes"
	"context"
	"github.com/bytedance/sonic"
	"github.com/jiu-u/oai-adapter/api"
	"io"
	"net/http"
)

func (c *Client) Embeddings(ctx context.Context, req *api.EmbeddingRequest) (io.ReadCloser, http.Header, error) {
	url := c.EndPoint + "v1beta/openai/embeddings"
	bodyBytes, err := sonic.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	return c.DoOAIJsonRequest(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
}

func (c *Client) EmbeddingsByBytes(ctx context.Context, req []byte) (io.ReadCloser, http.Header, error) {
	url := c.EndPoint + "v1beta/openai/embeddings"
	return c.DoOAIJsonRequest(ctx, "POST", url, bytes.NewBuffer(req))
}
