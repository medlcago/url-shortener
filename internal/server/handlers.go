package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/swaggo/http-swagger"
	"time"
	httpAuth "url-shortener/internal/auth/delivery/http"
	repoAuth "url-shortener/internal/auth/repository"
	serviceAuth "url-shortener/internal/auth/service"
	httpLinks "url-shortener/internal/links/delivery/http"
	repoLinks "url-shortener/internal/links/repository"
	serviceLinks "url-shortener/internal/links/service"
	"url-shortener/internal/middleware"
)

func (s *Server) MapHandlers() {
	// Init repositories
	linksRepo := repoLinks.NewLinksRepo(s.db)
	authRepo := repoAuth.NewAuthRepository(s.db)

	// Init services
	linksService := serviceLinks.NewLinkService(linksRepo, s.storage, s.logger)
	authService := serviceAuth.NewAuthService(authRepo, s.cfg, s.logger)

	// Init handlers
	linksHandlers := httpLinks.NewLinksHandlers(linksService)
	authHandlers := httpAuth.NewAuthHandlers(authService)

	// Init global middlewares
	mw := middleware.NewManager(authService, s.cfg, s.logger)
	s.app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		LimitReached: func(ctx fiber.Ctx) error {
			return fiber.ErrTooManyRequests
		},
	}))
	s.app.Use(requestid.New())
	s.app.Use(logger.New())
	s.app.Use(mw.DebugMiddleware())

	if s.cfg.Server.Mode == "Development" {
		s.app.Get("/swagger/*", adaptor.HTTPHandlerFunc(httpSwagger.Handler(
			httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
		)))
	}

	apiV1 := s.app.Group("/api/v1")

	// links endpoints
	{
		linksGroup := apiV1.Group("/links")
		httpLinks.MapLinksRoutes(linksGroup, linksHandlers, mw)
	}
	// auth endpoints
	{
		authGroup := apiV1.Group("/auth")
		httpAuth.MapAuthRoutes(authGroup, authHandlers, mw)
	}
}
