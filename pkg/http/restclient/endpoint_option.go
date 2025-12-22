package restclient

import (
	"fmt"
	"strings"
)

type EndpointOption func(r *request)

func WithUrlParam(key string, value interface{}) EndpointOption {
	return func(r *request) {
		r.URL = strings.ReplaceAll(r.URL, "{"+key+"}", fmt.Sprint(value))
	}
}

func WithQueryParam(key, value string) EndpointOption {
	return func(r *request) {
		r.SetQueryParam(key, value)
	}
}

func WithQueryParamList(key string, values []string) EndpointOption {
	return func(r *request) {
		r.SetQueryParam(key, strings.Join(values, ","))
	}
}

func WithQueryString(queryString string) EndpointOption {
	return func(r *request) {
		r.QueryString(queryString)
	}
}

func WithBody(body interface{}) EndpointOption {
	return func(r *request) {
		r.SetBody(body)
	}
}

func WithHeader(key string, value string) EndpointOption {
	return func(r *request) {
		r.SetHeader(key, value)
	}
}
func WithHeaders(hs map[string]string) EndpointOption {
	return func(r *request) {
		r.SetHeaders(hs)
	}
}

func WithBasicAuth(username, password string) EndpointOption {
	return func(r *request) {
		r.SetBasicAuth(username, password)
	}
}

func WithFailAt(failAt FailAtFunc) EndpointOption {
	return func(r *request) {
		r.SetFailAt(failAt)
	}
}
