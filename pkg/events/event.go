package events

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/jorgeAM/go-template/pkg/model"
)

var (
	ErrInvalidPayload  = errors.New("invalid payload")
	ErrInvalidReceiver = errors.New("receiver should be a pointer")
)

type Event struct {
	ID        string
	Topic     Topic
	Payload   interface{}
	Timestamp time.Time
	Version   int
}

func NewEvent(
	topic string,
	payload interface{},
) (*Event, error) {
	eventTopic, err := NewTopic(topic)
	if err != nil {
		return nil, err
	}

	if payload == nil {
		return nil, ErrInvalidPayload
	}

	return &Event{
		ID:        model.GenerateUUID().String(),
		Topic:     eventTopic,
		Payload:   payload,
		Timestamp: time.Now(),
		Version:   1,
	}, nil
}

func (e *Event) MarshalPayload() (json.RawMessage, error) {
	if b, ok := e.Payload.([]byte); ok {
		return b, nil
	}

	if b, ok := e.Payload.(json.RawMessage); ok {
		return b, nil
	}

	return json.Marshal(e.Payload)
}

func (e *Event) UnmarshalPayload(v interface{}) error {
	vValue := reflect.ValueOf(v)
	if vValue.Kind() != reflect.Ptr {
		return ErrInvalidReceiver
	}

	vValue = vValue.Elem()
	payloadValue := reflect.ValueOf(e.Payload)
	if vValue.Type() == payloadValue.Type() {
		vValue.Set(payloadValue)

		return nil
	}

	if b, ok := e.Payload.([]byte); ok {
		return json.Unmarshal(b, v)
	}

	if b, ok := e.Payload.(json.RawMessage); ok {
		return json.Unmarshal([]byte(b), v)
	}

	raw, err := e.MarshalPayload()
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, v)
}
