package eventbus

import (
	"context"

	"github.com/jorgeAM/go-template/internal/shared/events"
)

var _ Publisher = (*events.Collector)(nil)

//go:generate go tool mockgen -source=./publisher.go -destination=./mocks/publisher.go -package=mock -mock_names=Publisher=MockPublisher
type Publisher interface {
	Publish(ctx context.Context, events ...*events.Event) error
}
