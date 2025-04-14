package http

import (
	"github.com/gofiber/fiber/v3"
	"url-shortener/internal/links"
	"url-shortener/internal/middleware"
)

func MapLinksRoutes(linksGroup fiber.Router, handlers links.Handlers, mw *middleware.Manager) {
	linksGroup.Post("/", handlers.Create, mw.OptionalAuthMiddleware())
	linksGroup.Get("/:alias", handlers.Redirect)
}
