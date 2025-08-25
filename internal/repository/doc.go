package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paudarco/doc-storage/internal/entity"
	"github.com/paudarco/doc-storage/internal/errors"
)

type DocRepository struct {
	db *pgxpool.Pool
}

func NewDocRepository(db *pgxpool.Pool) *DocRepository {
	return &DocRepository{
		db: db,
	}
}

func (r *DocRepository) Create(ctx context.Context, doc *entity.Document) error {
	query := `INSERT INTO documents (id, user_id, name, is_file, public, mime, grant_list, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, doc.ID, doc.UserID, doc.Name, doc.IsFile, doc.Public, doc.Mime, doc.Grant, doc.CreatedAt)
	return err
}

func (r *DocRepository) GetByID(ctx context.Context, id string) (*entity.Document, error) {
	query := `SELECT d.id, d.user_id, d.name, d.is_file, d.public, d.mime, d.grant_list, d.created_at,
	                 u.login
	          FROM documents d
	          JOIN users u ON d.user_id = u.id
	          WHERE d.id = $1`
	doc := &entity.Document{}
	var ownerLogin string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&doc.ID, &doc.UserID, &doc.Name, &doc.IsFile, &doc.Public,
		&doc.Mime, &doc.Grant, &doc.CreatedAt, &ownerLogin,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrDocNotFound
		}
		return nil, err
	}
	return doc, nil
}

func (r *DocRepository) List(ctx context.Context, userID, loginFilter, keyFilter, valueFilter string, limit int) ([]*entity.Document, error) {
	targetUserID := userID
	if loginFilter != "" && loginFilter != userID {
		var userUUID string
		err := r.db.QueryRow(ctx, `SELECT id FROM users WHERE login = $1`, loginFilter).Scan(&userUUID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return []*entity.Document{}, nil
			}
			return nil, err
		}
		targetUserID = userUUID
	}

	baseQuery := `SELECT id, user_id, name, is_file, public, mime, grant_list, created_at FROM documents WHERE user_id = $1`
	args := []interface{}{targetUserID}
	argIndex := 2

	if keyFilter != "" && valueFilter != "" && keyFilter == "name" {
		baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+valueFilter+"%")
		argIndex++
	}

	baseQuery += " ORDER BY name ASC, created_at DESC"

	if limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
	}

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []*entity.Document
	for rows.Next() {
		doc := &entity.Document{}
		err := rows.Scan(&doc.ID, &doc.UserID, &doc.Name, &doc.IsFile, &doc.Public, &doc.Mime, &doc.Grant, &doc.CreatedAt)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return docs, nil
}

func (r *DocRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM documents WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.ErrDocNotFound
	}
	return nil
}
