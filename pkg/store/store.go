package store

import "context"

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error

	Scanner(ctx context.Context) (Scanner, error)
}

type Scanner interface {
	Next() bool
	Err() error

	Key() string
	Value() string

	Close() error
}
