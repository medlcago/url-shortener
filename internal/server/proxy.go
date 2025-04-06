package server

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
)

func (s *Server) MapProxy() {
	s.app.Get("/:alias", func(c fiber.Ctx) error {
		//url := fmt.Sprintf("http://localhost%s/api/v1/links/%s", s.cfg.Server.Port, c.Params("alias"))
		targetURL := fmt.Sprintf("/api/v1/links/%s", c.Params("alias"))
		if err := proxy.Do(c, targetURL); err != nil {
			return err
		}
		c.Response().Header.Del(fiber.HeaderServer)
		return nil
	})
}
