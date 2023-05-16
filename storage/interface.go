package storage

import "context"

type Storage interface {
	Save(ctx context.Context, key string, url string) error
	Load(ctx context.Context, key string) (string, error)
}
