package restclient

import (
	"context"
	"time"
)

type RetryConfig struct {
	Retries int
	Timeout time.Duration
}

type retry struct {
	config   RetryConfig
	endpoint Endpoint
}

func EndpointWithRetry(
	config RetryConfig,
	endpoint Endpoint,
) Endpoint {
	return &retry{
		config:   config,
		endpoint: endpoint,
	}
}

func (r *retry) DoRequest(ctx context.Context, opts ...EndpointOption) Response {
	var res Response

	for i := 0; i < r.config.Retries+1; i++ {
		res = r.endpoint.DoRequest(ctx, opts...)

		if res.Err() == nil {
			break
		}

		time.Sleep(r.config.Timeout)
	}

	return res
}

func (r *retry) Request() Request {
	return r.endpoint.Request()
}
