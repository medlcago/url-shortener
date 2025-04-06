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
}

func NewServer(app *fiber.App, cfg *config.Config, db *sqlx.DB, storage storage.Storage, logger logger.Logger) *Server {
	return &Server{
		app:     app,
		cfg:     cfg,
		db:      db,
		storage: storage,
		logger:  logger,
	}
}

func (s *Server) Run() error {
	s.MapHandlers()
	s.MapProxy()

	go func() {
		if err := s.app.Listen(s.cfg.Server.Port); err != nil {
			s.logger.Errorf("server listening error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.app.ShutdownWithContext(ctx)
}
