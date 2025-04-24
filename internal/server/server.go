package server

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/config"
	"url-shortener/pkg/logger"
	"url-shortener/pkg/storage"
)

const (
	ctxTimeout = 5
)

type Server struct {
	app     *fiber.App
	cfg     *config.Config
	db      *sqlx.DB
	storage storage.Storage
	logger  logger.Logger

	metrics *MetricsServer
}

func NewServer(app *fiber.App, cfg *config.Config, db *sqlx.DB, storage storage.Storage, logger logger.Logger) *Server {
	return &Server{
		app:     app,
		cfg:     cfg,
		db:      db,
		storage: storage,
		logger:  logger,
		metrics: NewMetricsServer(cfg.Server.MetricsPort),
	}
}

func (s *Server) Run() error {
	s.MapHandlers()

	go func() {
		if err := s.app.Listen(s.cfg.Server.Port); err != nil {
			s.logger.Errorf("API server listening error: %v", err)
		}
	}()

	go func() {
		if err := s.metrics.Listen(); err != nil {
			s.logger.Errorf("metrics server listening error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		s.logger.Errorf("API server shutdown error: %v", err)
		return err
	}

	if err := s.metrics.Shutdown(ctx); err != nil {
		s.logger.Errorf("metrics server shutdown error: %v", err)
		return err
	}

	s.logger.Info("Server Exited Properly")
	return nil
}
