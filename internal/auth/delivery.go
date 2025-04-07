package auth

import "github.com/gofiber/fiber/v3"

type Handlers interface {
	Login(ctx fiber.Ctx) error
	Register(ctx fiber.Ctx) error
	GetMe(ctx fiber.Ctx) error
}
