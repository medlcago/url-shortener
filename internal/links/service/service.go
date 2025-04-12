package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	"url-shortener/internal/links"
	"url-shortener/internal/links/repository"
	"url-shortener/internal/models"
	"url-shortener/pkg/http"
	"url-shortener/pkg/logger"
	"url-shortener/pkg/storage"
)

const (
	cacheKeyPrefix  = "shorturl:"
	defaultCacheTTL = 24 * time.Hour
)

func cacheKey(alias string) string {
	return cacheKeyPrefix + alias
}

type linkService struct {
	repo    links.Repository
	storage storage.Storage
	logger  logger.Logger
}

func NewLinkService(repo links.Repository, storage storage.Storage, logger logger.Logger) links.Service {
	return &linkService{repo: repo, storage: storage, logger: logger}
}

func (s *linkService) Create(ctx context.Context, link *models.Link) (*models.Link, error) {
	if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now().UTC()) {
		return nil, http.BadRequest
	}

	alias, err := randomString(8)
	if err != nil {
		s.logger.Errorf("failed to generate short URL: %v", err)
		return nil, http.InternalServerError
	}

	link.Alias = alias
	id, err := s.repo.Save(ctx, link)
	if err != nil {
		s.logger.Error("linksService.repo.Save:", err)
		return nil, http.InternalServerError
	}

	go func() {
		err := s.cacheLink(alias, link)
		if err != nil {
			s.logger.Warnf("failed to cache link (alias: %s): %v", alias, err)
		}
	}()

	s.logger.Infof("Created link with uid: %s", id)

	shortURL := fmt.Sprintf("%s/%s", link.BaseURL, alias)
	return models.NewLink(id, alias, link.OriginalURL, shortURL, link.ExpiresAt), nil
}

func (s *linkService) Resolve(ctx context.Context, alias string) (string, error) {
	if originalURL, err := s.storage.Get(ctx, cacheKey(alias)); err == nil {
		s.logger.Infof("Link found in cache (alias: %s)", alias)
		return originalURL.(string), nil
	}

	s.logger.Infof("Link not found in cache (alias: %s)", alias)

	link, err := s.repo.Get(ctx, alias)
	if err != nil {
		if errors.Is(err, repository.ErrLinkNotFound) {
			return "", http.NotFound
		}
		return "", http.InternalServerError
	}

	go func() {
		err := s.cacheLink(alias, link)
		if err != nil {
			s.logger.Warnf("failed to update cache for alias %s: %v", alias, err)
		}
	}()

	return link.OriginalURL, nil
}

func randomString(length int) (string, error) {
	bytes := make([]byte, length/2+1)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func (s *linkService) cacheLink(alias string, link *models.Link) error {
	var cacheTTL time.Duration
	if link.ExpiresAt == nil {
		cacheTTL = 30 * 24 * time.Hour
	} else {
		expiresIn := link.ExpiresAt.Sub(time.Now().UTC())
		cacheTTL = min(expiresIn, defaultCacheTTL)
		if cacheTTL <= 0 {
			return errors.New("link expired")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.storage.Set(ctx, cacheKey(alias), link.OriginalURL, cacheTTL); err != nil {
		return err
	}
	return nil
}
