package restclient

import (
	"context"
	"net/http"
)

type FailAtFunc func(req Request, res Response) error

func defaultFailAt(req Request, res Response) error {
	// 5xx status codes are not considered as errors.
	// It's the default behavior of Resty.
	return res.Err()
}

type Client interface {
	GET(url string, opts ...EndpointOption) Endpoint
	POST(url string, opts ...EndpointOption) Endpoint
	PUT(url string, opts ...EndpointOption) Endpoint
	PATCH(url string, opts ...EndpointOption) Endpoint
	DELETE(url string, opts ...EndpointOption) Endpoint
}

type Endpoint interface {
	DoRequest(ctx context.Context, opts ...EndpointOption) Response
	Request() Request
}

type Request interface {
	Do(ctx context.Context, opts ...EndpointOption) Response
	UrlParam(key string, value interface{}) Request
	QueryParam(key, value string) Request
	QueryParamList(key string, values []string) Request
	QueryString(queryString string) Request
	Body(body interface{}) Request
	Header(key string, value interface{}) Request
	Headers(hs map[string]string) Request
	BasicAuth(username, password string) Request
	SetFailAt(failAt FailAtFunc) Request
}

type Response interface {
	Body() []byte
	Status() string
	StatusCode() int
	Header() http.Header
	Err() error
}
