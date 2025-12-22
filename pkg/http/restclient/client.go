package restclient

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

var _ Client = (*restclient)(nil)

type restclient struct {
	client  *resty.Client
	baseUrl string
}

func New(cfg Config) *restclient {
	client := resty.New()
	client.SetBaseURL(cfg.BaseUrl)
	client.SetRetryCount(cfg.Retries)
	client.SetDebug(cfg.DebugMode)

	if cfg.CustomTransport != nil {
		client.SetTransport(cfg.CustomTransport)
	}

	if cfg.TimeoutMs != nil {
		client.SetTimeout(time.Duration(*cfg.TimeoutMs) * time.Millisecond)
	}

	if cfg.RetryWaitTimeMs != nil {
		retryWaitTime := *cfg.RetryWaitTimeMs
		if retryWaitTime <= 0 {
			retryWaitTime = 1
		}
		client.SetRetryWaitTime(time.Duration(retryWaitTime) * time.Millisecond)
	}

	rc := restclient{
		client:  client,
		baseUrl: cfg.BaseUrl,
	}

	return &rc
}

func (rc *restclient) GET(urlFormat string, opts ...EndpointOption) Endpoint {
	var endpoint Endpoint = &endpoint{
		client:    rc.client,
		urlFormat: urlFormat,
		method:    http.MethodGet,
		headers:   make(http.Header),
		opts:      opts,
	}

	return endpoint
}

func (rc *restclient) POST(urlFormat string, opts ...EndpointOption) Endpoint {
	var endpoint Endpoint = &endpoint{
		client:    rc.client,
		urlFormat: urlFormat,
		method:    http.MethodPost,
		headers:   make(http.Header),
		opts:      opts,
	}

	return endpoint
}

func (rc *restclient) PUT(urlFormat string, opts ...EndpointOption) Endpoint {
	var endpoint Endpoint = &endpoint{
		client:    rc.client,
		urlFormat: urlFormat,
		method:    http.MethodPut,
		headers:   make(http.Header),
		opts:      opts,
	}

	return endpoint
}

func (rc *restclient) PATCH(urlFormat string, opts ...EndpointOption) Endpoint {
	var endpoint Endpoint = &endpoint{
		client:    rc.client,
		urlFormat: urlFormat,
		method:    http.MethodPatch,
		headers:   make(http.Header),
		opts:      opts,
	}

	return endpoint
}

func (rc *restclient) DELETE(urlFormat string, opts ...EndpointOption) Endpoint {
	var endpoint Endpoint = &endpoint{
		client:    rc.client,
		urlFormat: urlFormat,
		method:    http.MethodDelete,
		headers:   make(http.Header),
		opts:      opts,
	}

	return endpoint
}
