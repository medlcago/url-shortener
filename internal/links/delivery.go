package links

import "github.com/gofiber/fiber/v3"

type Handlers interface {
	Create(ctx fiber.Ctx) error
	Redirect(ctx fiber.Ctx) error
	GetAll(ctx fiber.Ctx) error
}
