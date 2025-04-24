package service

import (
	"context"
	"errors"
	"time"
	"url-shortener/internal/models"
)

const (
	cacheKeyPrefix  = "shorturl:"
	defaultCacheTTL = 24 * time.Hour
)

func cacheKey(alias string) string {
	return cacheKeyPrefix + alias
}

func (s *linkService) cacheLink(ctx context.Context, alias string, link *models.Link) error {
	var cacheTTL time.Duration
	if link.ExpiresAt == nil {
		cacheTTL = 30 * 24 * time.Hour
	} else {
		ttl := link.ExpiresAt.Sub(time.Now().UTC())
		cacheTTL = min(ttl, defaultCacheTTL)
		if cacheTTL <= 0 {
			return errors.New("link expired")
		}
	}

	if err := s.storage.Set(ctx, cacheKey(alias), link.OriginalURL, cacheTTL); err != nil {
		return err
	}
	return nil
}
