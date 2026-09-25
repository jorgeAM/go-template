package criteria

import (
	"errors"
	"strings"
)

var ErrInvalidOrderType = errors.New("invalid order type")

type OrderType string

const (
	ASC  OrderType = "ASC"
	DESC OrderType = "DESC"
)

var allowedOrderType = map[string]OrderType{
	ASC.String():  ASC,
	DESC.String(): DESC,
}

func NewOrderType(o string) (OrderType, error) {
	if orderType, ok := allowedOrderType[strings.ToUpper(o)]; ok {
		return orderType, nil
	}

	return "", ErrInvalidOrderType
}

func (o OrderType) String() string {
	return string(o)
}
