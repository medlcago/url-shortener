package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/sirupsen/logrus"
	"net/http/httputil"
)

func (m *Manager) DebugMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		if !m.cfg.Server.Debug {
			return ctx.Next()
		}

		httpReq, err := adaptor.ConvertRequest(ctx, false)
		if err != nil {
			m.logger.WithFields(logrus.Fields{
				"op":         "adaptor.ConvertRequest",
				"middleware": "DebugMiddleware",
				"error":      err.Error(),
			}).Error("failed to convert request")

			return err
		}

		dump, err := httputil.DumpRequest(httpReq, true)
		if err != nil {
			m.logger.WithFields(logrus.Fields{
				"op":         "httputil.DumpRequest",
				"middleware": "DebugMiddleware",
				"error":      err.Error(),
			}).Error("failed to dump request")

			return err
		}

		m.logger.WithFields(logrus.Fields{
			"op":         "httputil.DumpRequest",
			"middleware": "DebugMiddleware",
			"request":    string(dump),
		}).Info("dump request")

		return ctx.Next()
	}
}
