package middleware

import (
	"url-shortener/config"
	"url-shortener/pkg/logger"
)

type Manager struct {
	cfg    *config.Config
	logger logger.Logger
}

func NewManager(cfg *config.Config, logger logger.Logger) *Manager {
	return &Manager{
		cfg:    cfg,
		logger: logger,
	}
}
