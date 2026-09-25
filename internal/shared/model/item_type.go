package model

import (
	"errors"
	"strings"
)

var (
	ErrInvalidItemType = errors.New("invalid item type")
)

type ItemType string

const (
	Bag      ItemType = "BAG"
	Suitcase ItemType = "SUITCASE"
	Cap      ItemType = "CAP"
	Bike     ItemType = "BIKE"
	Car      ItemType = "CAR"
	Other    ItemType = "OTHER"
)

var allowedItemType = map[string]ItemType{
	Bag.String():      Bag,
	Suitcase.String(): Suitcase,
	Cap.String():      Cap,
	Bike.String():     Bike,
	Car.String():      Car,
	Other.String():    Other,
}

func NewItemType(t string) (ItemType, error) {
	if status, ok := allowedItemType[strings.ToUpper(t)]; ok {
		return status, nil
	}

	return "", ErrInvalidItemType
}

func (r ItemType) String() string {
	return string(r)
}
