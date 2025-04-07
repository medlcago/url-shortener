package http

import (
	"github.com/gofiber/fiber/v3"
	"url-shortener/internal/auth"
)

func MapAuthRoutes(authGroup fiber.Router, handlers auth.Handlers) {
	authGroup.Post("/login", handlers.Login)
}
