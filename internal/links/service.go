package links

import (
	"context"
	"url-shortener/internal/models"
)

type Service interface {
	Create(ctx context.Context, link *models.Link) (*models.Link, error)
	Resolve(ctx context.Context, alias string) (string, error)
}
