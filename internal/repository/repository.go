package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paudarco/doc-storage/internal/entity"
)

type User interface {
	Create(ctx context.Context, user *entity.User) error
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
}

type Doc interface {
	Create(ctx context.Context, doc *entity.Document) error
	GetByID(ctx context.Context, id string) (*entity.Document, error)
	List(ctx context.Context, userID, loginFilter, keyFilter, valueFilter string, limit int) ([]*entity.Document, error)
	Delete(ctx context.Context, id string) error
}

type Repository struct {
	User
	Doc
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		User: NewUserRepository(db),
		Doc:  NewDocRepository(db),
	}
}
