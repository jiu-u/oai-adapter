package openai

import (
	"errors"
	"io"
	"net/http"
)

func NoImplementMethod() (io.ReadCloser, http.Header, error) {
	return nil, nil, errors.New("not implement")
}
