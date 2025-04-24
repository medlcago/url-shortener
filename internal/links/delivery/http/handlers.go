package http

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"time"
	"url-shortener/internal/auth"
	"url-shortener/internal/links"
	"url-shortener/internal/models"
	"url-shortener/pkg/http"
	"url-shortener/pkg/pagination"
)

const (
	defaultCtxTimeout = 5 * time.Second
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
//
//	@Summary		Create a short URL
//	@Description	Generates a short alias for the provided URL
//	@Tags			links
//	@Accept			json
//	@Produce		json
//	@Param			request		body		request						true	"URL data for shortening"
//	@Param			X-Anonymous	header		boolean						true	"Создать ссылку без привязки к пользователю"	default(false)
//	@Success		201			{object}	http.Response[models.Link]	"URL successfully shortened"
//	@Failure		400			{object}	http.Response[any]			"Invalid input data"
//	@Failure		401			{object}	http.Response[any]			"Unauthorized - Missing or invalid token"
//	@Failure		500			{object}	http.Response[any]			"Internal server error"
//	@Security		BearerAuth
//	@Router			/links [post]
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
		if authData, exists := ctx.Locals("authData").(*auth.Data); exists {
			link.OwnerID = &authData.User.ID
		} else {
			return http.InvalidCredentials
		}
	}

	link.BaseURL = ctx.BaseURL()

	c, cancel := context.WithTimeout(ctx.Context(), defaultCtxTimeout)
	defer cancel()

	res, err := h.linksService.Create(c, &link)
	if err != nil {
		return err
	}

	data := http.OK[*models.Link](res)
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

// Redirect godoc
//
//	@Summary		Redirect by short URL
//	@Description	Performs redirect from short URL to original
//	@Tags			links
//	@Param			alias	path	string	true	"Short URL identifier"
//	@Success		308		"Permanent redirect to original URL"
//	@Failure		404		{object}	http.Response[any]	"URL not found"
//	@Failure		500		{object}	http.Response[any]	"Internal server error"
//	@Router			/links/{alias} [get]
func (h *linksHandlers) Redirect(ctx fiber.Ctx) error {
	alias := fiber.Params[string](ctx, "alias")
	c, cancel := context.WithTimeout(ctx.Context(), defaultCtxTimeout)
	defer cancel()

	url, err := h.linksService.Resolve(c, alias)
	if err != nil {
		return err
	}
	return ctx.Redirect().Status(fiber.StatusPermanentRedirect).To(url)
}

// GetAll godoc
//
//	@Summary		Get all user's links
//	@Description	Retrieves all shortened URLs for the authenticated user with pagination
//	@Tags			links
//	@Security		BearerAuth
//	@Param			limit	query		int								false	"Number of items per page"	default(10)	minimum(1)	maximum(100)
//	@Param			offset	query		int								false	"Offset for pagination"		default(0)	minimum(0)
//	@Success		200		{object}	http.Response[[]models.Link]	"List of user's links with pagination info"
//	@Failure		401		{object}	http.Response[any]				"Unauthorized - Missing or invalid token"
//	@Failure		500		{object}	http.Response[any]				"Internal server error"
//	@Router			/links [get]
func (h *linksHandlers) GetAll(ctx fiber.Ctx) error {
	authData := fiber.Locals[*auth.Data](ctx, "authData")
	p := pagination.FromContext(ctx)

	res, total, err := h.linksService.GetAll(ctx.Context(), ctx.BaseURL(), authData.User.ID, p.Limit, p.Offset)
	if err != nil {
		return err
	}

	data := http.OK[[]models.Link](res, http.MetaData{"total": total})
	return ctx.JSON(data)
}
