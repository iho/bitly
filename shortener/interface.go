package shortener

import "context"

type Shortener interface {
	Save(ctx context.Context, url string) (string, error)
	Load(ctx context.Context, key string) (string, error)
}
