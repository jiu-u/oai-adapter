package base

import (
	v1 "github.com/jiu-u/oai-adapter/api/v1"
	"io"
	"net/http"
)

func NoImplementMethod() (io.ReadCloser, http.Header, error) {
	return nil, nil, v1.NoImplementError
}
