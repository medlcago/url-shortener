package http

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"url-shortener/internal/auth"
	"url-shortener/internal/models"
	"url-shortener/pkg/http"
)

type authHandlers struct {
	authService auth.Service
}

func NewAuthHandlers(authService auth.Service) auth.Handlers {
	return &authHandlers{authService: authService}
}

func (h *authHandlers) Login(ctx fiber.Ctx) error {
	var user models.User
	if err := ctx.Bind().JSON(&user); err != nil {
		return err
	}

	token, err := h.authService.Login(context.Background(), &user)
	if err != nil {
		return err
	}

	data := http.NewResponse[*auth.Token](true, "", token)
	return ctx.JSON(data)
}

func (h *authHandlers) Register(ctx fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotImplemented)
}

func (h *authHandlers) GetMe(ctx fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotImplemented)
}
