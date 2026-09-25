package eventbus

import (
	"context"

	"github.com/jorgeAM/go-template/internal/shared/events"
)

type Handler interface {
	HandlerID() string
	Handle(ctx context.Context, event *events.Event) error
}
