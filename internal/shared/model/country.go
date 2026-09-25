package model

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCountry             = errors.New("invalid_country")
	ErrCountryDoesNotHaveCurrency = errors.New("country_does_not_have_currency")
)

type Country string

const (
	PE Country = "PE"
	AR Country = "AR"
	US Country = "US"
)

var allowedCountry = map[string]Country{
	PE.String(): PE,
	AR.String(): AR,
	US.String(): US,
}

func NewCountry(t string) (Country, error) {
	if country, ok := allowedCountry[strings.ToUpper(t)]; ok {
		return country, nil
	}

	return "", ErrInvalidCountry
}

func (c Country) String() string {
	return string(c)
}

var currencyByCountry = map[Country]Currency{
	PE: PEN,
	AR: ARS,
	US: USD,
}

func (c Country) GetCurrency() (Currency, error) {
	if currency, ok := currencyByCountry[c]; ok {
		return currency, nil
	}

	return "", ErrCountryDoesNotHaveCurrency
}
