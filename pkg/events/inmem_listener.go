package events

import (
	"context"

	"github.com/jorgeAM/go-template/pkg/log"
)

var _ Listener = (*InMemoryListener)(nil)

type InMemoryListener struct {
	bus      *InMemoryEventBus
	handlers map[Topic]Handler

	// events without handler
	missed []*Event
}

func NewInMemoryListener(handlers map[Topic]Handler) *InMemoryListener {
	return &InMemoryListener{
		bus:      getInMemoryEventBus(),
		handlers: handlers,
		missed:   []*Event{},
	}
}

func (i *InMemoryListener) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			for _, event := range i.bus.Drain() {
				handler, ok := i.handlers[event.Topic]
				if !ok {
					log.Warn(ctx, "event don't have handler", log.WithString("topic", event.Topic.String()), log.WithObject("event", event))
					i.missed = append(i.missed, event)
					continue
				}

				if err := handler.Handle(ctx, event); err != nil {
					log.Error(ctx, "error handling event from inMemListener", log.WithString("topic", event.Topic.String()), log.WithError(err))
					continue
				}
			}
		}
	}
}
