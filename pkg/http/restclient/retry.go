package restclient

import (
	"context"
	"time"
)

type RetryConfig struct {
	// Retries to be executed after first try.
	// If retry = 0, at least one request will be executed.
	// If retry > 0, the first request will be executed without considering
	// this value and after that retries will be considered.
	Retries int
	// Sleep time between requests
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
