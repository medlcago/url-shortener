package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
	"url-shortener/config"
	"url-shortener/internal/auth"
	"url-shortener/internal/auth/repository"
	"url-shortener/internal/models"
	"url-shortener/pkg/http"
	"url-shortener/pkg/jwt"
	"url-shortener/pkg/logger"
	"url-shortener/pkg/storage"
)

type authService struct {
	repo    auth.Repository
	storage storage.Storage
	cfg     *config.Config
	logger  logger.Logger
}

func NewAuthService(repo auth.Repository, storage storage.Storage, cfg *config.Config, logger logger.Logger) auth.Service {
	return &authService{repo: repo, storage: storage, cfg: cfg, logger: logger}
}

func (s *authService) Login(ctx context.Context, user *models.User) (*auth.Token, error) {
	u, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, repository.UserNotFound) {
			s.logger.WithFields(logrus.Fields{
				"email":    user.Email,
				"endpoint": "Login",
				"type":     "auth",
			}).Warn("login attempt with non-existent email")

			return nil, http.InvalidCredentials
		}

		s.logger.WithFields(logrus.Fields{
			"email":    user.Email,
			"error":    err.Error(),
			"endpoint": "Login",
			"op":       "repo.GetByEmail",
		}).Error("failed to get user by email")

		return nil, http.InternalServerError
	}

	if err = u.ComparePasswords(user.Password); err != nil {
		return nil, http.InvalidCredentials
	}

	tokenPair, err := jwt.NewTokenPair(u.ID.String(), s.cfg.Server.JwtSecretKey)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id":  u.ID,
			"op":       "jwt.NewTokenPair",
			"endpoint": "Login",
			"error":    err.Error(),
		}).Error("failed to create token pair")

		return nil, http.InternalServerError
	}

	if err = s.repo.UpdateLastLogin(ctx, u.ID, time.Now().UTC()); err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id":  u.ID,
			"op":       "repo.UpdateLastLogin",
			"endpoint": "Login",
			"error":    err.Error(),
		}).Error("failed to update last login time")
	}

	return auth.NewToken(tokenPair.AccessToken, tokenPair.RefreshToken), nil
}

func (s *authService) Register(ctx context.Context, user *models.User) (*auth.Token, error) {
	if _, err := s.repo.GetByEmail(ctx, user.Email); err == nil {
		return nil, http.ExistsEmailError
	}

	if err := user.PrepareCreate(); err != nil {
		s.logger.WithFields(logrus.Fields{
			"email":    user.Email,
			"op":       "user.PrepareCreate",
			"endpoint": "Register",
			"error":    err.Error(),
		}).Error("failed to prepare user for registration")

		return nil, http.InternalServerError
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.WithFields(logrus.Fields{
			"email":    user.Email,
			"op":       "repo.Create",
			"endpoint": "Register",
			"error":    err.Error(),
		}).Error("failed to create user")

		return nil, http.InternalServerError
	}

	tokenPair, err := jwt.NewTokenPair(user.ID.String(), s.cfg.Server.JwtSecretKey)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id":  user.ID,
			"op":       "jwt.NewTokenPair",
			"endpoint": "Register",
			"error":    err.Error(),
		}).Error("failed to create token pair")

		return nil, http.InternalServerError
	}

	return auth.NewToken(tokenPair.AccessToken, tokenPair.RefreshToken), nil
}

func (s *authService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.UserNotFound) {
			return nil, http.NotFound
		}

		s.logger.WithFields(logrus.Fields{
			"user_id":  id,
			"op":       "repo.GetByID",
			"endpoint": "GetByID",
			"error":    err.Error(),
		}).Error("failed to get user by ID")

		return nil, http.InternalServerError
	}

	return user, nil
}

func (s *authService) RefreshToken(ctx context.Context, data *auth.Data) (*auth.Token, error) {
	key := s.refreshTokenKey(data.Token)
	exists, err := s.storage.Exists(ctx, key)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id":  data.User.ID,
			"op":       "storage.Exists",
			"endpoint": "RefreshToken",
			"error":    err.Error(),
		}).Error("failed to check token existence")

		return nil, http.InternalServerError
	}

	if exists {
		return nil, http.InvalidToken
	}

	if err := s.storage.Set(ctx, s.refreshTokenKey(data.Token), data.Token, data.TTL); err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id":  data.User.ID,
			"op":       "storage.Set",
			"endpoint": "RefreshToken",
			"error":    err.Error(),
		}).Error("failed to set refresh token")

		return nil, http.InternalServerError
	}

	res, err := jwt.NewTokenPair(data.User.ID.String(), s.cfg.Server.JwtSecretKey)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id":  data.User.ID,
			"op":       "jwt.NewTokenPair",
			"endpoint": "RefreshToken",
			"error":    err.Error(),
		}).Error("failed to create new token pair")

		return nil, http.InternalServerError
	}

	return auth.NewToken(res.AccessToken, res.RefreshToken), nil
}

func (s *authService) refreshTokenKey(token string) string {
	return fmt.Sprintf("refreshToken:%s", token)
}
