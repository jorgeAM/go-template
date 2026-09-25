package restclient

import (
	"net/http"
)

type Config struct {
	BaseUrl         string            `json:"base_url"`
	TimeoutMs       *int              `json:"timeout_ms"`
	Retries         int               `json:"retries"`
	CustomTransport http.RoundTripper `json:"-"`
	RetryWaitTimeMs *int              `json:"retry_wait_time_ms"`
	DebugMode       bool              `json:"debug_mode"`
}
