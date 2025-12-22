package events

import "context"

var _ Publisher = (*InMemoryPublisher)(nil)

type InMemoryPublisher struct {
	bus *InMemoryEventBus
}

func NewInMemoryPublisher() *InMemoryPublisher {
	return &InMemoryPublisher{bus: getInMemoryEventBus()}
}

func (p *InMemoryPublisher) Publish(ctx context.Context, events ...*Event) error {
	p.bus.Add(events...)
	return nil
}
