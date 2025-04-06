package http

import (
	"github.com/gofiber/fiber/v3"
	"url-shortener/internal/links"
)

func MapLinksRoutes(linksGroup fiber.Router, handlers links.Handlers) {
	linksGroup.Post("/", handlers.Create)
	linksGroup.Get("/:alias", handlers.Redirect).Name("LinkRedirect")
}
