package criteria

import "errors"

var ErrInvalidPageSize = errors.New("invalid page size")

type PageSize int

func NewPageSize(value int) (PageSize, error) {
	if value <= 0 || value > 200 {
		return 0, ErrInvalidPageSize
	}

	return PageSize(value), nil
}

func (p PageSize) Raw() int {
	return int(p)
}
