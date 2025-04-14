package models

import (
	"github.com/google/uuid"
	"time"
)

type Link struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Alias       string     `json:"alias" db:"alias"`
	OriginalURL string     `json:"original_url" validate:"http_url" db:"original_url"`
	ShortURL    string     `json:"short_url"`
	BaseURL     string     `json:"-"`
	OwnerID     *uuid.UUID `json:"owner_id" db:"owner_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at" db:"expires_at"`
}
