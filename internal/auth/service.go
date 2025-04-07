package auth

import (
	"context"
	"github.com/google/uuid"
	"url-shortener/internal/models"
)

type Token struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}

type Service interface {
	Login(ctx context.Context, user *models.User) (*Token, error)
	Register(ctx context.Context, user *models.User) (*Token, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Token, error)
}
