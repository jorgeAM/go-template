package events

import "sync"

var (
	inMemoryBus     *InMemoryEventBus
	inMemoryBusOnce sync.Once
)

type InMemoryEventBus struct {
	mu     sync.Mutex
	events []*Event
}

func getInMemoryEventBus() *InMemoryEventBus {
	inMemoryBusOnce.Do(func() {
		inMemoryBus = &InMemoryEventBus{
			events: []*Event{},
		}
	})

	return inMemoryBus
}

func (i *InMemoryEventBus) Add(events ...*Event) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.events = append(i.events, events...)
}

func (i *InMemoryEventBus) Drain() []*Event {
	i.mu.Lock()
	defer i.mu.Unlock()

	events := i.events
	i.events = make([]*Event, 0)

	return events
}
