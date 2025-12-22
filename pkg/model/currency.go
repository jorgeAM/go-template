package model

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCurrency = errors.New("invalid currency")
)

type Currency string

const (
	PEN Currency = "PEN"
	ARS Currency = "ARS"
	USD Currency = "USD"
)

var allowedCurrency = map[string]Currency{
	PEN.String(): PEN,
	ARS.String(): ARS,
	USD.String(): USD,
}

func NewCurrency(t string) (Currency, error) {
	if currency, ok := allowedCurrency[strings.ToUpper(t)]; ok {
		return currency, nil
	}

	return "", ErrInvalidCountry
}

func (r Currency) String() string {
	return string(r)
}
