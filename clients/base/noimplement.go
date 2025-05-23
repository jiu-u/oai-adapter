package base

import (
	"context"
	"github.com/jiu-u/oai-adapter/common"
	"io"
	"net/http"
)

func NoImplementMethod(_ context.Context, _ any) (io.ReadCloser, http.Header, error) {
	return nil, nil, common.NoImplementError
}
