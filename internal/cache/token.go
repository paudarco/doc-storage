package cache

import (
	"context"
	"time"

	"github.com/paudarco/doc-storage/internal/errors"
	"github.com/redis/go-redis/v9"
)

type TokenCache struct {
	cache *redis.Client
	exp   time.Duration
}

func NewTokenCache(cache *redis.Client, exp time.Duration) *TokenCache {
	return &TokenCache{
		cache: cache,
		exp:   exp,
	}
}

func (c *TokenCache) SetToken(ctx context.Context, token, userID string) error {
	key := TokenPrefix + token
	return c.cache.Set(ctx, key, userID, c.exp).Err()
}

func (c *TokenCache) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	key := TokenPrefix + token
	userID, err := c.cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.ErrWrongToken
	} else if err != nil {
		return "", err
	}
	return userID, nil
}

func (c *TokenCache) DeleteToken(ctx context.Context, token string) error {
	key := TokenPrefix + token
	return c.cache.Del(ctx, key).Err()
}
