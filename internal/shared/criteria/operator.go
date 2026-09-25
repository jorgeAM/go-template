package criteria

import (
	"errors"
	"strings"
)

var ErrInvalidOperator = errors.New("invalid operator")

type Operator string

const (
	EQUAL Operator = "EQ"
	GT    Operator = "GT"
	LT    Operator = "LT"
)

var allowedOperator = map[string]Operator{
	EQUAL.String(): EQUAL,
	GT.String():    GT,
	LT.String():    LT,
}

func NewOperator(o string) (Operator, error) {
	if operator, ok := allowedOperator[strings.ToUpper(o)]; ok {
		return operator, nil
	}

	return "", ErrInvalidOperator
}

func (o Operator) String() string {
	return string(o)
}
