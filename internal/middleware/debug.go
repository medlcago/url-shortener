package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"net/http/httputil"
)

func (m *Manager) DebugMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		if m.cfg.Server.Debug {
			httpReq, err := adaptor.ConvertRequest(ctx, false)
			if err != nil {
				return err
			}
			dump, err := httputil.DumpRequest(httpReq, true)
			if err != nil {
				return err
			}
			m.logger.Infof("\nRequest dump begin :--------------\n\n%s\n\nRequest dump end :--------------", dump)
		}
		return ctx.Next()
	}
}
