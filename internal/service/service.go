package service

import (
	"context"
	"encoding/json"

	"github.com/paudarco/doc-storage/internal/cache"
	"github.com/paudarco/doc-storage/internal/config"
	"github.com/paudarco/doc-storage/internal/entity"
	"github.com/paudarco/doc-storage/internal/repository"
	"github.com/sirupsen/logrus"
)

type Auth interface {
}

type User interface {
	Register(ctx context.Context, login, password string) error
	Authenticate(ctx context.Context, login, password string) (string, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	InvalidateToken(ctx context.Context, token string) error
}

type Doc interface {
	Create(ctx context.Context, userID string, meta map[string]interface{}, jsonData json.RawMessage, fileData []byte) (*entity.Document, error)
	List(ctx context.Context, userID, loginFilter, keyFilter, valueFilter string, limit int) ([]*entity.Document, error)
	GetByID(ctx context.Context, userID, docID string) (*entity.Document, error)
	checkAccess(ctx context.Context, doc *entity.Document, userID string) error
	Delete(ctx context.Context, userID, docID string) error
}

type Service struct {
	Auth
	User
	Doc
}

func NewService(repo *repository.Repository, cache *cache.Cache, cfg *config.Config, log *logrus.Logger) *Service {
	return &Service{
		User: NewUserService(repo.User, cache.Token, cfg),
		Doc:  NewDocService(repo.Doc, repo.User, cache.Doc, log),
	}
}
