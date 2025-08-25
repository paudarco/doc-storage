package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/paudarco/doc-storage/internal/cache"
	"github.com/paudarco/doc-storage/internal/entity"
	"github.com/paudarco/doc-storage/internal/errors"
	"github.com/paudarco/doc-storage/internal/repository"
	"github.com/sirupsen/logrus"
)

type DocService struct {
	docRepo  repository.Doc
	userRepo repository.User
	cache    cache.Doc
	log      *logrus.Logger
}

func NewDocService(docRepo repository.Doc, userRepo repository.User, cache cache.Doc, log *logrus.Logger) *DocService {
	return &DocService{
		docRepo:  docRepo,
		userRepo: userRepo,
		cache:    cache,
		log:      log,
	}
}

func (s *DocService) Create(ctx context.Context, userID string, meta map[string]interface{}, jsonData json.RawMessage, fileData []byte) (*entity.Document, error) {
	doc := &entity.Document{
		ID:        uuid.New().String(),
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if name, ok := meta["name"].(string); ok {
		doc.Name = name
	} else {
		return nil, errors.ErrMetaNameRequired
	}

	if isFile, ok := meta["file"].(bool); ok {
		doc.IsFile = isFile
	} else {
		doc.IsFile = false
	}

	if public, ok := meta["public"].(bool); ok {
		doc.Public = public
	} else {
		doc.Public = false
	}

	if mime, ok := meta["mime"].(string); ok && doc.IsFile {
		doc.Mime = mime
	}

	if grantList, ok := meta["grant"].([]interface{}); ok {
		for _, g := range grantList {
			if login, ok := g.(string); ok {
				doc.Grant = append(doc.Grant, login)
			}
		}
	}

	doc.JSONData = jsonData
	doc.FileData = fileData

	err := s.docRepo.Create(ctx, doc)
	if err != nil {
		s.log.Errorf("failed to create document in DB: %v", err)
		return nil, fmt.Errorf("failed to create document in DB")
	}

	_ = s.cache.InvalidateUserDocLists(ctx, userID)

	return doc, nil
}

func (s *DocService) List(ctx context.Context, userID, loginFilter, keyFilter, valueFilter string, limit int) ([]*entity.Document, error) {
	targetUserID := userID
	if loginFilter != "" {
		var userUUID uuid.UUID
		_, err := s.userRepo.GetByID(ctx, loginFilter)
		if err != nil {
			user, err := s.userRepo.GetByLogin(ctx, loginFilter)
			if err != nil {
				return []*entity.Document{}, nil
			}
			userUUID = user.ID
		} else {
			userUUID = uuid.MustParse(loginFilter)
		}
		targetUserID = userUUID.String()
	}

	cacheKey := cache.BuildDocListCacheKey(userID, loginFilter, keyFilter, valueFilter, limit)

	cachedData, err := s.cache.GetDocList(ctx, cacheKey)
	if err != nil {
		s.log.Printf("Error getting doc list from cache: %v", err)
	} else if cachedData != nil {
		var docs []*entity.Document
		if err := json.Unmarshal(*cachedData, &docs); err != nil {
			s.log.Printf("Error unmarshalling cached doc list: %v", err)
		} else {
			return docs, nil
		}
	}

	docs, err := s.docRepo.List(ctx, targetUserID, loginFilter, keyFilter, valueFilter, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get document list from DB: %w", err)
	}

	dataToCache, err := json.Marshal(docs)
	if err != nil {
		s.log.Errorf("Error marshalling doc list for cache: %v", err)
	} else {
		_ = s.cache.SetDocList(ctx, cacheKey, dataToCache)
	}

	return docs, nil
}

func (s *DocService) GetByID(ctx context.Context, userID, docID string) (*entity.Document, error) {
	cachedData, err := s.cache.GetDoc(ctx, docID)
	if err != nil {

		s.log.Errorf("Error getting doc from cache: %v", err)
	} else if cachedData != nil {
		var doc entity.Document
		if err := json.Unmarshal(*cachedData, &doc); err != nil {
			s.log.Printf("Error unmarshalling cached doc: %v", err)
		} else {
			if accessErr := s.checkAccess(ctx, &doc, userID); accessErr != nil {
				return nil, accessErr
			}

			return &doc, nil
		}
	}

	doc, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return nil, err
	}

	if accessErr := s.checkAccess(ctx, doc, userID); accessErr != nil {
		return nil, accessErr
	}

	dataToCache, err := json.Marshal(doc)
	if err != nil {
		s.log.Errorf("Error marshalling doc for cache: %v", err)
	} else {
		_ = s.cache.SetDoc(ctx, docID, dataToCache)
	}

	return doc, nil
}

func (s *DocService) checkAccess(ctx context.Context, doc *entity.Document, userID string) error {
	if doc.UserID != userID {
		if !doc.Public {
			hasAccess := false
			for _, grantedLogin := range doc.Grant {
				currentUser, err := s.userRepo.GetByID(ctx, userID)
				if err != nil {
					return err
				}
				if currentUser.Login == grantedLogin {
					hasAccess = true
					break
				}
			}
			if !hasAccess {
				return errors.ErrAccessDenied
			}
		}
	}
	return nil
}

func (s *DocService) Delete(ctx context.Context, userID, docID string) error {

	doc, err := s.docRepo.GetByID(ctx, docID)
	if err != nil {
		return err
	}

	if doc.UserID != userID {
		return errors.ErrAccessDenied
	}

	err = s.docRepo.Delete(ctx, docID)
	if err != nil {
		return err
	}

	_ = s.cache.DeleteDoc(ctx, docID)
	_ = s.cache.InvalidateUserDocLists(ctx, doc.UserID)

	return nil
}
