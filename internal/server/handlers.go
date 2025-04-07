package server

import (
	"github.com/gofiber/fiber/v3/middleware/requestid"
	httpAuth "url-shortener/internal/auth/delivery/http"
	repoAuth "url-shortener/internal/auth/repository"
	serviceAuth "url-shortener/internal/auth/service"
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
	authRepo := repoAuth.NewAuthRepository(s.db)

	// Init services
	linksService := serviceLinks.NewLinkService(linksRepo, s.storage, s.logger)
	authService := serviceAuth.NewAuthService(authRepo, s.cfg, s.logger)

	// Init handlers
	linksHandlers := httpLinks.NewLinksHandlers(linksService)
	authHandlers := httpAuth.NewAuthHandlers(authService)

	apiV1 := s.app.Group("/api/v1")

	linksGroup := apiV1.Group("/links")
	authGroup := apiV1.Group("/auth")

	httpLinks.MapLinksRoutes(linksGroup, linksHandlers)
	httpAuth.MapAuthRoutes(authGroup, authHandlers)
}
