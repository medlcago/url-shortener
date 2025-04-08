package http

import (
	"github.com/gofiber/fiber/v3"
	"url-shortener/internal/auth"
	"url-shortener/internal/middleware"
)

func MapAuthRoutes(authGroup fiber.Router, handlers auth.Handlers, mw *middleware.Manager) {
	authGroup.Post("/login", handlers.Login)
	authGroup.Post("/register", handlers.Register)
	authGroup.Get("/me", handlers.GetMe, mw.CurrentUserMiddleware())
}
