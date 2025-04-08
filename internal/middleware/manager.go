package middleware

import (
	"url-shortener/config"
	"url-shortener/internal/auth"
	"url-shortener/pkg/logger"
)

type Manager struct {
	authService auth.Service
	cfg         *config.Config
	logger      logger.Logger
}

func NewManager(authService auth.Service, cfg *config.Config, logger logger.Logger) *Manager {
	return &Manager{
		authService: authService,
		cfg:         cfg,
		logger:      logger,
	}
}
