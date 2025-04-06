package server

import (
	"github.com/gofiber/fiber/v3/middleware/requestid"
	httpLinks "url-shortener/internal/links/delivery/http"
	repoLinks "url-shortener/internal/links/repository"
	serviceLinks "url-shortener/internal/links/service"
	"url-shortener/internal/middleware"
)

func (s *Server) MapHandlers() {
	mw := middleware.NewManager(s.cfg, s.logger)
	s.app.Use(requestid.New())
	s.app.Use(mw.DebugMiddleware())

	// Init repositories
	linksRepo := repoLinks.NewLinksRepo(s.db)

	// Init services
	linksService := serviceLinks.NewLinkService(linksRepo, s.storage, s.logger)

	// Init handlers
	linksHandlers := httpLinks.NewLinksHandlers(linksService)

	apiV1 := s.app.Group("/api/v1")

	linksGroup := apiV1.Group("/links")

	httpLinks.MapLinksRoutes(linksGroup, linksHandlers)
}
