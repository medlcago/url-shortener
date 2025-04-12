package middleware

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"url-shortener/pkg/http"
	"url-shortener/pkg/utils"
)

func extractJWTClaims(ctx fiber.Ctx, secret string) (*utils.Claims, error) {
	req, err := adaptor.ConvertRequest(ctx, false)
	if err != nil {
		return nil, err
	}

	claims, err := utils.ExtractClaimsFromRequest(req, secret)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (m *Manager) AuthJWTMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		claims, err := extractJWTClaims(ctx, m.cfg.Server.JwtSecretKey)
		if err != nil {
			return http.InvalidToken
		}

		ctx.Locals("user_id", claims.ID)
		return ctx.Next()
	}
}

func (m *Manager) CurrentUserMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		claims, err := extractJWTClaims(ctx, m.cfg.Server.JwtSecretKey)
		if err != nil {
			return http.InvalidToken
		}

		user, err := m.authService.GetByID(context.Background(), claims.ID)
		if err != nil {
			return http.InvalidCredentials
		}

		user.SanitizePassword()

		ctx.Locals("user", user)
		return ctx.Next()
	}
}
