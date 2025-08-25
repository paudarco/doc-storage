package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/paudarco/doc-storage/internal/errors"
)

type Client struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewClient(secret, issuer string, ttl int) *Client {
	return &Client{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    time.Duration(ttl) * time.Hour,
	}
}

func (c *Client) GenerateAccessToken(userID string) (string, error) {
	now := time.Now().UTC()
	claims := &jwt.RegisteredClaims{
		Issuer:    c.issuer,
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(now.Add(c.ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        uuid.NewString(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.secret)
}

func (c *Client) ValidateAccessToken(tokenString string) (string, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return c.secret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := tok.Claims.(*jwt.RegisteredClaims)
	if !ok || !tok.Valid {
		return "", errors.ErrInvalidToken
	}

	// Проверяем только время жизни (exp)
	if claims.ExpiresAt == nil || time.Now().UTC().After(claims.ExpiresAt.Time) {
		return "", errors.ErrTokenExpired
	}

	userID := claims.Subject
	if userID == "" {
		return "", errors.ErrInvalidToken
	}
	return userID, nil
}
