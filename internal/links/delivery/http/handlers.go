package http

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"url-shortener/internal/links"
	"url-shortener/internal/models"
	"url-shortener/pkg/http"
)

type linksHandlers struct {
	linksService links.Service
}

func NewLinksHandlers(linksService links.Service) links.Handlers {
	return &linksHandlers{linksService: linksService}
}

func (h *linksHandlers) Create(ctx fiber.Ctx) error {
	var link models.Link
	if err := ctx.Bind().JSON(&link); err != nil {
		return err
	}

	link.BaseURL = ctx.BaseURL()
	res, err := h.linksService.Create(context.Background(), &link)
	if err != nil {
		return err
	}

	data := http.NewResponse[*models.Link](true, "", res)
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

func (h *linksHandlers) Redirect(ctx fiber.Ctx) error {
	alias := fiber.Params[string](ctx, "alias")
	url, err := h.linksService.Resolve(context.Background(), alias)
	if err != nil {
		return err
	}
	return ctx.Redirect().Status(fiber.StatusPermanentRedirect).To(url)
}
