package pagination

import (
	"github.com/gofiber/fiber/v3"
)

const (
	DefaultLimit    = 10
	DefaultOffset   = 0
	DefaultMaxLimit = 100
)

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func FromContext(ctx fiber.Ctx) *Pagination {
	p := &Pagination{
		Limit:  DefaultLimit,
		Offset: DefaultOffset,
	}

	if err := ctx.Bind().Query(p); err != nil {
		return p
	}

	if p.Limit < 1 || p.Limit > DefaultMaxLimit {
		p.Limit = DefaultLimit
	}

	if p.Offset < 1 {
		p.Offset = DefaultOffset
	}

	return p
}
