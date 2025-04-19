package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/big"
	"time"
	"url-shortener/internal/links"
	"url-shortener/internal/links/repository"
	"url-shortener/internal/models"
	"url-shortener/pkg/http"
	"url-shortener/pkg/logger"
	"url-shortener/pkg/storage"
)

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

	alias, err := s.generateRandomString(8)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"op":       "linkService.generateRandomString",
			"endpoint": "Create",
			"error":    err.Error(),
		}).Error("failed to generate short URL")

		return nil, http.InternalServerError
	}

	link.Alias = alias

	if err := s.repo.Save(ctx, link); err != nil {
		s.logger.WithFields(logrus.Fields{
			"op":       "repo.Save",
			"endpoint": "Create",
			"error":    err.Error(),
		}).Error("failed to save short URL")

		return nil, http.InternalServerError
	}

	go func() {
		err := s.cacheLink(alias, link)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"op":       "linkService.cacheLink",
				"endpoint": "Create",
				"alias":    alias,
				"error":    err.Error(),
			}).Warnf("failed to cache link")
		}
	}()

	s.logger.WithFields(logrus.Fields{
		"op":       "linkService.Create",
		"endpoint": "Create",
		"alias":    alias,
	}).Info("link created")

	shortURL := fmt.Sprintf("%s/%s", link.BaseURL, alias)
	link.ShortURL = shortURL
	return link, nil
}

func (s *linkService) Resolve(ctx context.Context, alias string) (string, error) {
	if originalURL, err := s.storage.Get(ctx, cacheKey(alias)); err == nil {
		s.logger.WithFields(logrus.Fields{
			"op":       "linkService.Resolve",
			"endpoint": "Resolve",
			"alias":    alias,
			"cached":   true,
		}).Info("link found in cache")

		return originalURL.(string), nil
	}

	s.logger.WithFields(logrus.Fields{
		"op":       "linkService.Resolve",
		"endpoint": "Resolve",
		"alias":    alias,
		"cached":   false,
	}).Info("link not found in cache")

	link, err := s.repo.Get(ctx, alias)
	if err != nil {
		if errors.Is(err, repository.ErrLinkNotFound) {
			return "", http.NotFound
		}

		s.logger.WithFields(logrus.Fields{
			"op":       "repo.Get",
			"endpoint": "Resolve",
			"error":    err.Error(),
		}).Error("failed to get link")

		return "", http.InternalServerError
	}

	go func() {
		err := s.cacheLink(alias, link)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"op":       "cacheLink",
				"endpoint": "Resolve",
				"alias":    alias,
				"error":    err.Error(),
			}).Warnf("failed to update cache link")
		}
	}()

	return link.OriginalURL, nil
}

func (s *linkService) generateRandomString(length int) (string, error) {
	if length < 1 {
		return "", errors.New("length must be > 0")
	}

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
