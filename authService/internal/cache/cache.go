package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrorTokenNotSet     = errors.New("token is not set")
	ErrorTokenAlreadySet = errors.New("token already set")
)

type ICache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Ping(ctx context.Context) error
	Close() error
}
