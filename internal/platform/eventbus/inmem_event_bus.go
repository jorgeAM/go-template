package eventbus

import (
	"sync"

	"github.com/jorgeAM/go-template/internal/shared/events"
)

var (
	inMemoryBus     *InMemoryEventBus
	inMemoryBusOnce sync.Once
)

type InMemoryEventBus struct {
	mu     sync.Mutex
	events []*events.Event
}

func getInMemoryEventBus() *InMemoryEventBus {
	inMemoryBusOnce.Do(func() {
		inMemoryBus = &InMemoryEventBus{
			events: []*events.Event{},
		}
	})

	return inMemoryBus
}

func (i *InMemoryEventBus) Add(events ...*events.Event) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.events = append(i.events, events...)
}

func (i *InMemoryEventBus) Drain() []*events.Event {
	i.mu.Lock()
	defer i.mu.Unlock()

	drained := i.events
	i.events = make([]*events.Event, 0)

	return drained
}
