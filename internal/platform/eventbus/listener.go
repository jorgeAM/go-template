package eventbus

import "context"

type Listener interface {
	Listen(ctx context.Context)
}
