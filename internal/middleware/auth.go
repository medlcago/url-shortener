package middleware

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/google/uuid"
	"time"
	"url-shortener/internal/auth"
	"url-shortener/pkg/http"
	"url-shortener/pkg/jwt"
)

func extractJWTClaims(ctx fiber.Ctx, secret string) (string, *jwt.Claims, error) {
	req, err := adaptor.ConvertRequest(ctx, false)
	if err != nil {
		return "", nil, err
	}

	token, err := jwt.ExtractTokenFromRequest(req)
	if err != nil {
		return "", nil, err
	}

	claims, err := jwt.ExtractClaimsFromToken(token, secret)
	if err != nil {
		return "", nil, err
	}

	return token, claims, nil
}

func (m *Manager) AuthJWTMiddleware(tokenType jwt.TokenType) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		token, claims, err := extractJWTClaims(ctx, m.cfg.Server.JwtSecretKey)
		if err != nil {
			return http.InvalidToken
		}

		if claims.Type != tokenType {
			return http.InvalidToken
		}

		userID, err := uuid.Parse(claims.ID)
		if err != nil {
			m.logger.Errorf("failed to parse user id from claims: %s", err.Error())
			return http.InvalidToken
		}

		user, err := m.authService.GetByID(context.Background(), userID)
		if err != nil {
			m.logger.Errorf("failed to get user by id: %s", err.Error())
			return http.InvalidToken
		}

		user.SanitizePassword()

		authData := &auth.Data{Token: token, User: user, TTL: time.Until(claims.ExpiresAt.Time)}
		ctx.Locals("authData", authData)
		return ctx.Next()
	}
}

func (m *Manager) OptionalAuthMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		anonymousHeader := fiber.GetReqHeader[bool](ctx, "X-Anonymous", false)
		ctx.Locals("anonymous", anonymousHeader)
		if anonymousHeader {
			return ctx.Next()
		}

		return m.AuthJWTMiddleware(jwt.Access)(ctx)
	}
}
