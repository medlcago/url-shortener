package links

import (
	"context"
	"github.com/google/uuid"
	"url-shortener/internal/models"
)

type Service interface {
	Create(ctx context.Context, link *models.Link) (*models.Link, error)
	Resolve(ctx context.Context, alias string) (string, error)
	GetAll(ctx context.Context, baseURL string, ownerID uuid.UUID, limit, offset int) ([]models.Link, int64, error)
}
