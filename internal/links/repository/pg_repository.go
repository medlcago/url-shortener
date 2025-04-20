package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
	"url-shortener/internal/links"
	"url-shortener/internal/models"
)

var (
	ErrLinkNotFound = errors.New("link not found")
)

type linksRepo struct {
	db *sqlx.DB
}

func NewLinksRepo(db *sqlx.DB) links.Repository {
	return &linksRepo{db: db}
}

func (r *linksRepo) Save(ctx context.Context, link *models.Link) error {
	if err := r.db.GetContext(ctx, link, createLink, link.Alias, link.OriginalURL, link.ExpiresAt, link.OwnerID); err != nil {
		return err
	}
	return nil
}

func (r *linksRepo) Get(ctx context.Context, alias string) (*models.Link, error) {
	var link models.Link
	err := r.db.GetContext(ctx, &link, getLinkByAlias, alias, time.Now().UTC())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLinkNotFound
	}
	return &link, err
}

func (r *linksRepo) Exists(ctx context.Context, alias string) (bool, error) {
	var exists bool
	if err := r.db.GetContext(ctx, &exists, existsLink, alias); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *linksRepo) SelectAll(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]models.Link, int64, error) {
	var data []models.Link
	var count int64

	if err := r.db.SelectContext(ctx, &data, getAllUserLinks, ownerID, limit, offset); err != nil {
		return nil, 0, err
	}

	if err := r.db.GetContext(ctx, &count, countUserLinks, ownerID); err != nil {
		return nil, 0, err
	}

	return data, count, nil
}
