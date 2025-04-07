package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"url-shortener/pkg/http"
	"url-shortener/pkg/utils"
)

func (m *Manager) AuthJWTMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		req, err := adaptor.ConvertRequest(ctx, false)
		if err != nil {
			m.logger.Errorf("error converting request: %v", err)
			return err
		}

		claims, err := utils.ExtractClaimsFromRequest(req, m.cfg.Server.JwtSecretKey)
		if err != nil {
			m.logger.Errorf("error extracting claims: %v", err)
			return http.WithMessage(err.Error()).SetStatus(fiber.StatusUnauthorized)
		}

		ctx.Locals("user_id", claims.ID)
		return ctx.Next()
	}
}
