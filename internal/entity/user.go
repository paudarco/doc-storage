package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Login     string    `json:"login" db:"login"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type RegisterRequest struct {
	AdminToken string `json:"token" binding:"required"`
	Login      string `json:"login" binding:"required"`
	Pswd       string `json:"pswd" binding:"required"`
}

type LoginRequest struct {
	Login string `json:"login" binding:"required"`
	Pswd  string `json:"pswd" binding:"required"`
}
