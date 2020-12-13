package rpcRouter

import "context"

type Flusher interface {
	Flush(ctx context.Context) error
}
