package events

import "context"

type Handler interface {
	HandlerID() string
	Handle(ctx context.Context, event *Event) error
}
