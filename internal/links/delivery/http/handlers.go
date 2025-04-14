package http

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"time"
	"url-shortener/internal/links"
	"url-shortener/internal/models"
	"url-shortener/pkg/http"
)

type request struct {
	OriginalURL string     `json:"original_url" validate:"http_url" example:"https://github.com/"`
	ExpiresAt   *time.Time `json:"expires_at" example:"2006-01-02T15:04:05Z"`
}

type linksHandlers struct {
	linksService links.Service
}

func NewLinksHandlers(linksService links.Service) links.Handlers {
	return &linksHandlers{linksService: linksService}
}

// Create godoc
// @Summary Create a short URL
// @Description Generates a short alias for the provided URL
// @Tags links
// @Accept json
// @Produce json
// @Param request body request true "URL data for shortening"
// @Param   X-Anonymous header boolean true "Создать ссылку без привязки к пользователю" default(false)
// @Success 201 {object} http.Response[models.Link] "URL successfully shortened"
// @Failure 400 {object} http.Response[any] "Invalid input data"
// @Failure 401 {object} http.Response[any] "Unauthorized - Missing or invalid token"
// @Failure 500 {object} http.Response[any] "Internal server error"
// @Security BearerAuth
// @Router /links [post]
func (h *linksHandlers) Create(ctx fiber.Ctx) error {
	var req request
	if err := ctx.Bind().JSON(&req); err != nil {
		return err
	}

	link := models.Link{
		OriginalURL: req.OriginalURL,
		ExpiresAt:   req.ExpiresAt,
	}
	isAnonymous := fiber.Locals[bool](ctx, "anonymous")
	if !isAnonymous {
		if user, exists := ctx.Locals("user").(*models.User); exists {
			link.OwnerID = &user.ID
		} else {
			return fiber.ErrUnauthorized
		}
	}

	link.BaseURL = ctx.BaseURL()
	res, err := h.linksService.Create(context.Background(), &link)
	if err != nil {
		return err
	}

	data := http.OK[*models.Link](res)
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

// Redirect godoc
// @Summary Redirect by short URL
// @Description Performs redirect to the original URL using short alias
// @Tags links
// @Param alias path string true "Short URL identifier"
// @Success 308 "Permanent redirect to original URL"
// @Failure 404 {object} http.Response[any] "URL not found"
// @Failure 500 {object} http.Response[any] "Internal server error"
// @Router /{alias} [get]
func (h *linksHandlers) Redirect(ctx fiber.Ctx) error {
	alias := fiber.Params[string](ctx, "alias")
	url, err := h.linksService.Resolve(context.Background(), alias)
	if err != nil {
		return err
	}
	return ctx.Redirect().Status(fiber.StatusPermanentRedirect).To(url)
}
