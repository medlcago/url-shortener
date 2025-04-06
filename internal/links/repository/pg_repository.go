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

func (r *linksRepo) Save(ctx context.Context, link *models.Link) (uuid.UUID, error) {
	res := r.db.QueryRowxContext(ctx, createLink, link.Alias, link.OriginalURL, link.ExpiresAt)
	var uid uuid.UUID
	if err := res.Scan(&uid); err != nil {
		return uuid.Nil, err
	}
	return uid, nil
}

func (r *linksRepo) Get(ctx context.Context, alias string) (*models.Link, error) {
	var link models.Link
	err := r.db.QueryRowContext(ctx, getLinkByAlias, alias, time.Now().UTC()).Scan(
		&link.ID, &link.OriginalURL, &link.Alias, &link.ExpiresAt, &link.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLinkNotFound
	}
	return &link, err
}

func (r *linksRepo) Exists(ctx context.Context, alias string) (bool, error) {
	var exists bool
	if err := r.db.QueryRowxContext(ctx, existsLink, alias).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
