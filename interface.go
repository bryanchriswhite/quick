package quick

import "context"

type IncrementServer interface {
	Start(ctx context.Context) error
	IsClosed() bool
	Close() error
}

type IncrementClient interface {
	Increment(amount uint64) (newCount uint64, err error)
	Close() error
}
