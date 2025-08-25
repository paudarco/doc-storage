package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/paudarco/doc-storage/internal/errors"
	"github.com/redis/go-redis/v9"
)

type DocCache struct {
	cache *redis.Client
	exp   time.Duration
}

func NewDocCache(cache *redis.Client, exp time.Duration) *DocCache {
	return &DocCache{
		cache: cache,
		exp:   exp,
	}
}

func (c *DocCache) SetDoc(ctx context.Context, id string, docData []byte) error {
	key := DocPrefix + id
	return c.cache.Set(ctx, key, docData, c.exp).Err()
}

func (c *DocCache) GetDoc(ctx context.Context, id string) (*[]byte, error) {
	key := DocPrefix + id
	val, err := c.cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.ErrDocNotFound
	} else if err != nil {
		return nil, err
	}

	data := []byte(val)
	return &data, nil
}

func (c *DocCache) DeleteDoc(ctx context.Context, id string) error {
	key := DocPrefix + id
	return c.cache.Del(ctx, key).Err()
}

func (c *DocCache) SetDocList(ctx context.Context, cacheKey string, listData []byte) error {
	return c.cache.Set(ctx, cacheKey, listData, c.exp).Err()
}

func (c *DocCache) GetDocList(ctx context.Context, cacheKey string) (*[]byte, error) {
	val, err := c.cache.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, errors.ErrDocListNotFound
	} else if err != nil {
		return nil, err
	}

	data := []byte(val)
	return &data, nil
}

func BuildDocListCacheKey(userID, loginFilter, keyFilter, valueFilter string, limit int) string {
	// Для простоты используем форматирование строки. В production лучше использовать хеширование.
	return fmt.Sprintf("%s:%s:%s:%s:%s:%d", DocListPrefix, userID, loginFilter, keyFilter, valueFilter, limit)
}

// InvalidateUserDocLists Инвалидирует все списки документов конкретного пользователя
// Это делается путем установки флага или удаления ключей по паттерну.
// Более эффективный способ: использовать Redis Sets или Tags.
// Здесь реализация через паттерн SCAN.
func (c *DocCache) InvalidateUserDocLists(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s:%s:*", DocListPrefix, userID)
	iter := c.cache.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := c.cache.Del(ctx, iter.Val()).Err()
		if err != nil {
			// Логируем ошибку, но продолжаем
			log.Printf("Error deleting cache key %s: %v", iter.Val(), err)
		}
	}
	return iter.Err()
}

// InvalidateAllDocListsForUserLogin Инвалидирует списки, где мог быть документ этого пользователя
// Например, если документ стал публичным или изменился grant.
// Это более сложная операция, требует отслеживания всех возможных списков.
// Для упрощения можно инвалидировать все списки, связанные с этим пользователем как владельцем.
// Или использовать более сложную систему тегов/индексов в Redis.
// Пока реализуем как инвалидацию по userID владельца.
func (c *DocCache) InvalidateAllDocListsForUserLogin(ctx context.Context, userLogin string) error {
	// Эта операция сложнее, так как мы не знаем ID пользователя по логину напрямую в кэше.
	// В реальном приложении либо храним маппинг login->userID в кэше, либо инвалидируем иначе.
	// Пока оставим заглушку или реализуем позже.
	// Например, можно хранить в Redis SET "user_login:{login}" -> userID
	// Или просто инвалидировать все списки (неэффективно).
	// Для демонстрации: инвалидируем по userID, который мы должны знать из контекста создания/обновления.
	// В service.DocService.Create/Update/Delete мы будем вызывать InvalidateUserDocLists(userID)
	// Это покроет большинство случаев.
	return nil
}
