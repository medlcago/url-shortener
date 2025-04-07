package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
	"url-shortener/config"
	"url-shortener/internal/auth"
	"url-shortener/internal/auth/repository"
	"url-shortener/internal/models"
	"url-shortener/pkg/http"
	"url-shortener/pkg/logger"
	"url-shortener/pkg/utils"
)

type authService struct {
	repo   auth.Repository
	cfg    *config.Config
	logger logger.Logger
}

func NewAuthService(repo auth.Repository, cfg *config.Config, logger logger.Logger) auth.Service {
	return &authService{repo: repo, cfg: cfg, logger: logger}
}

func (s *authService) Login(ctx context.Context, user *models.User) (*auth.Token, error) {
	u, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, repository.UserNotFound) {
			return nil, http.WithMessage("invalid credentials").SetStatus(401)
		}
		s.logger.WithFields(logrus.Fields{
			"email": user.Email,
			"error": err.Error(),
		}).Error("failed to get user")

		return nil, http.WithMessage("failed to login")
	}

	if err = u.ComparePasswords(user.Password); err != nil {
		return nil, http.WithMessage("invalid credentials").SetStatus(401)
	}

	token, err := utils.GenerateJWTToken(u, s.cfg.Server.JwtSecretKey)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": u.ID,
		}).Error("failed to generate JWT")

		return nil, http.WithMessage("login temporarily unavailable")
	}

	if err = s.repo.UpdateLastLogin(ctx, u.ID, time.Now().UTC()); err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": u.ID,
		}).Error("failed to update last login time")
	}

	return &auth.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}

func (s *authService) Register(ctx context.Context, user *models.User) (*auth.Token, error) {
	return nil, nil
}

func (s *authService) GetByID(ctx context.Context, id uuid.UUID) (*auth.Token, error) {
	return nil, nil
}
