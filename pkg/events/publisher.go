package events

import "context"

//go:generate mockgen -source=./publisher.go -destination=./mocks/publisher.go -package=mock -mock_names=Publisher=MockPublisher
type Publisher interface {
	Publish(ctx context.Context, events ...*Event) error
}
