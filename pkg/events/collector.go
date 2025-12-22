package events

import (
	"context"
	"sync"
)

var _ Publisher = (*Collector)(nil)

type Collector struct {
	mux    sync.RWMutex
	events []*Event
}

func NewCollector() *Collector {
	return &Collector{
		events: make([]*Event, 0),
	}
}

func (c *Collector) Publish(ctx context.Context, events ...*Event) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.events = append(c.events, events...)

	return nil
}

func (c *Collector) Collect(events ...*Event) error {
	return c.Publish(context.Background(), events...)
}

func (c *Collector) Events() []*Event {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.events
}

func (c *Collector) Drain() []*Event {
	c.mux.Lock()
	defer c.mux.Unlock()

	events := c.events
	c.events = make([]*Event, 0)

	return events
}
