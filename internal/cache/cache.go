package cache

import (
	"context"
	"time"

	"github.com/paudarco/doc-storage/internal/config"
	"github.com/redis/go-redis/v9"
)

const (
	TokenPrefix    = "token:"
	DocPrefix      = "doc:"
	DocListPrefix  = "doc_list:"
	UserDocsPrefix = "user_docs:"
)

type Token interface {
	SetToken(ctx context.Context, token, userID string) error
	GetUserIDByToken(ctx context.Context, token string) (string, error)
	DeleteToken(ctx context.Context, token string) error
}

type Doc interface {
	SetDoc(ctx context.Context, id string, docData []byte) error
	GetDoc(ctx context.Context, id string) (*[]byte, error)
	DeleteDoc(ctx context.Context, id string) error
	SetDocList(ctx context.Context, cacheKey string, listData []byte) error
	GetDocList(ctx context.Context, cacheKey string) (*[]byte, error)
	InvalidateUserDocLists(ctx context.Context, userID string) error
}

type Cache struct {
	Token
	Doc
}

func NewCache(cache *redis.Client, cfg *config.Config) *Cache {
	return &Cache{
		Token: NewTokenCache(cache, time.Duration(cfg.AccessTTL)*time.Hour),
		Doc:   NewDocCache(cache, time.Duration(cfg.DocTTL)*time.Hour),
	}
}
