package auth

import (
	"context"
	"github.com/google/uuid"
	"time"
	"url-shortener/internal/models"
)

type Token struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}

func NewToken(accessToken, refreshToken string) *Token {
	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}
}

type Data struct {
	Token string        `json:"token"`
	TTL   time.Duration `json:"ttl"`
	User  *models.User  `json:"user"`
}

type Service interface {
	Login(ctx context.Context, user *models.User) (*Token, error)
	Register(ctx context.Context, user *models.User) (*Token, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	RefreshToken(ctx context.Context, data *Data) (*Token, error)
}
