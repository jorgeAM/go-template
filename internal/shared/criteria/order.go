package criteria

import "errors"

var ErrInvalidOrderBy = errors.New("invalid order by")

type Order struct {
	OrderBy   string
	OrderType OrderType
}

func NewOrder(orderBy, orderType string) (*Order, error) {
	if orderBy == "" {
		return nil, ErrInvalidOrderBy
	}

	oType, err := NewOrderType(orderType)
	if err != nil {
		return nil, err
	}

	return &Order{
		OrderBy:   orderBy,
		OrderType: oType,
	}, nil
}
