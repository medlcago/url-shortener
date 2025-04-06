package links

import (
	"context"
	"github.com/google/uuid"
	"url-shortener/internal/models"
)

type Repository interface {
	Save(ctx context.Context, link *models.Link) (uuid.UUID, error)
	Get(ctx context.Context, alias string) (*models.Link, error)
	Exists(ctx context.Context, alias string) (bool, error)
}
