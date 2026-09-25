package eventbus

import (
	"context"
	"github.com/jorgeAM/go-template/internal/shared/events"

	"github.com/jorgeAM/go-template/internal/platform/log"
)

var _ Listener = (*InMemoryListener)(nil)

type InMemoryListener struct {
	bus       *InMemoryEventBus
	handlers  map[events.Topic]Handler
	unhandled []*events.Event
}

func NewInMemoryListener(handlers map[events.Topic]Handler) *InMemoryListener {
	return &InMemoryListener{
		bus:       getInMemoryEventBus(),
		handlers:  handlers,
		unhandled: []*events.Event{},
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
					i.unhandled = append(i.unhandled, event)
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
