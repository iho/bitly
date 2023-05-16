package shortener

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/iho/bitly/storage"
	"go.uber.org/fx"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

var r *rand.Rand

func init() {
	seed := time.Now().UnixNano()
	r = rand.New(rand.NewSource(seed))
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

type shortener struct {
	storage storage.Storage
}

func New(storage storage.Storage) Shortener {
	return &shortener{storage: storage}
}

func (s *shortener) Save(ctx context.Context, url string) (string, error) {
	key := ""
	// looking for a free short key
	for i := 0; i < 10; i++ {
		key = RandStringBytes(8)
		url, _ := s.storage.Load(ctx, key)
		if url == "" { // key not found
			break
		}
	}

	if key == "" {
		return "", errors.New("can't find free key in a reasonable time")
	}

	return key, s.storage.Save(ctx, key, url)
}

func (s *shortener) Load(ctx context.Context, key string) (string, error) {
	return s.storage.Load(ctx, key)
}

// Module for go fx
var Module = fx.Options(fx.Provide(New))
