package restclient

import "github.com/go-resty/resty/v2"

var _ Response = (*response)(nil)

type response struct {
	err error
	*resty.Response
}

func (r *response) Err() error {
	return r.err
}
