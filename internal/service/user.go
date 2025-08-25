package service

import (
	"context"
	"fmt"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/paudarco/doc-storage/internal/cache"
	"github.com/paudarco/doc-storage/internal/config"
	"github.com/paudarco/doc-storage/internal/entity"
	"github.com/paudarco/doc-storage/internal/errors"
	"github.com/paudarco/doc-storage/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.User
	cfg      *config.Config
	cache    cache.Token
}

func NewUserService(userRepo repository.User, cache cache.Token, cfg *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		cache:    cache,
		cfg:      cfg,
	}
}

func (s *UserService) Register(ctx context.Context, login, password string) error {
	if err := validateLogin(login); err != nil {
		return err
	}

	if err := validatePassword(password); err != nil {
		return err
	}

	_, err := s.userRepo.GetByLogin(ctx, login)
	if err == nil {
		return errors.ErrUserAlreadyExist
	} else if err != errors.ErrUserNotFound {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		ID:        uuid.New(),
		Login:     login,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Authenticate(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return "", errors.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.ErrInvalidCredentials
	}

	token := uuid.New().String()

	err = s.cache.SetToken(ctx, token, user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) ValidateToken(ctx context.Context, token string) (string, error) {
	userID, err := s.cache.GetUserIDByToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("failed to get user ID from cache: %w", err)
	}
	if userID == "" {
		return "", fmt.Errorf("invalid or expired token")
	}
	return userID, nil
}

func (s *UserService) InvalidateToken(ctx context.Context, token string) error {
	return s.cache.DeleteToken(ctx, token)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func validateLogin(login string) error {
	if len(login) < 8 {
		return errors.ErrWrongLoginLength
	}

	var (
		lowerCount, upperCount, digitCount int
	)
	for _, c := range login {
		switch {
		case c >= 'a' && c <= 'z':
			lowerCount++
		case c >= 'A' && c <= 'Z':
			upperCount++
		case c >= '0' && c <= '9':
			digitCount++
		default:
			return errors.ErrLoginWithoutLatin
		}
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var (
		lowerCount, upperCount, digitCount, symbolCount int
	)
	for _, c := range password {
		switch {
		case c >= 'a' && c <= 'z':
			lowerCount++
		case c >= 'A' && c <= 'Z':
			upperCount++
		case c >= '0' && c <= '9':
			digitCount++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			symbolCount++
		default:
			return errors.ErrPswrdWithoutLatin
		}
	}

	if lowerCount < 2 {
		return errors.ErrPswrdWithoutLower
	}
	if upperCount < 2 {
		return errors.ErrPswrdWithoutUpper
	}
	if digitCount < 1 {
		return errors.ErrPswrdWithoutDigit
	}
	if symbolCount < 1 {
		return errors.ErrPswrdWithoutSymbol
	}

	return nil
}
