package events

import "context"

//go:generate go tool mockgen -source=./publisher.go -destination=./mocks/publisher.go -package=mock -mock_names=Publisher=MockPublisher
type Publisher interface {
	Publish(ctx context.Context, events ...*Event) error
}
