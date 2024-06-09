package quick

import "context"

type ServerI interface {
	Start(ctx context.Context) error
	IsClosed() bool
	Close() error
}

type ClientI interface {
	Increment(amount uint64) (newCount uint64, err error)
	Close() error
}
