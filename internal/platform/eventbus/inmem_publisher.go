package eventbus

import (
	"context"

	"github.com/jorgeAM/go-template/internal/shared/events"
)

var _ Publisher = (*InMemoryPublisher)(nil)

type InMemoryPublisher struct {
	bus *InMemoryEventBus
}

func NewInMemoryPublisher() *InMemoryPublisher {
	return &InMemoryPublisher{bus: getInMemoryEventBus()}
}

func (p *InMemoryPublisher) Publish(ctx context.Context, events ...*events.Event) error {
	p.bus.Add(events...)
	return nil
}
