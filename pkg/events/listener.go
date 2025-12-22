package events

import "context"

type Listener interface {
	Listen(ctx context.Context)
}
