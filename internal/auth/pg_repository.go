package auth

import (
	"context"
	"github.com/google/uuid"
	"time"
	"url-shortener/internal/models"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error
}
