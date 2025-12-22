package criteria

import (
	"errors"
)

var ErrInvalidPage = errors.New("invalid page")

type Page int

func NewPage(value int) (Page, error) {
	if value <= 0 {
		return 0, ErrInvalidPage
	}

	return Page(value), nil
}

func (p Page) Raw() int {
	return int(p)
}
