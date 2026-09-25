package events

import "errors"

var (
	ErrInvalidTopic = errors.New("invalid topic")
)

type Topic string

func NewTopic(topic string) (Topic, error) {
	if topic == "" {
		return "", ErrInvalidTopic
	}

	return Topic(topic), nil
}

func (t Topic) String() string {
	return string(t)
}
