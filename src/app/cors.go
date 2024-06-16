package app

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/ochom/gutils/env"
)

func b2cCors() fiber.Handler {
	origins := env.Get("B2C_ALLOWED_ORIGINS")
	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return strings.Contains(origins, origin)
		},
	})
}
