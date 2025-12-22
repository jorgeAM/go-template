package restclient

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var _ Endpoint = (*endpoint)(nil)

type endpoint struct {
	urlFormat string
	method    string
	headers   http.Header
	client    *resty.Client
	opts      []EndpointOption
}

func (e *endpoint) DoRequest(ctx context.Context, opts ...EndpointOption) Response {
	return newRequest(e).Do(ctx, opts...)
}

func (e *endpoint) Request() Request {
	return newRequest(e)
}
