package app

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/ochom/mpesa/src/app/config"
)

func b2cCors() fiber.Handler {
	crs := cors.ConfigDefault
	crs.AllowOriginsFunc = func(origin string) bool {
		return strings.Contains(config.B2CAllowedOrigins, origin)
	}

	return cors.New(crs)
}

func taxCors() fiber.Handler {
	crs := cors.ConfigDefault
	crs.AllowOriginsFunc = func(origin string) bool {
		return strings.Contains(config.TaxAllowedOrigins, origin)
	}

	return cors.New(crs)
}
