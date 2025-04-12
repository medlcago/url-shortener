package models

import (
	"github.com/google/uuid"
	"time"
)

type Link struct {
	ID          uuid.UUID `json:"id"`
	Alias       string    `json:"alias"`
	OriginalURL string    `json:"original_url" validate:"http_url"`
	ShortURL    string    `json:"short_url"`
	BaseURL     string    `json:"-"`
	// OwnerID     *uuid.UUID `json:"owner_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func NewLink(id uuid.UUID, alias string, originalURL string, shortURL string, expiresAt *time.Time) *Link {
	return &Link{
		ID:          id,
		Alias:       alias,
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		CreatedAt:   time.Now().UTC(),
		ExpiresAt:   expiresAt,
	}
}
