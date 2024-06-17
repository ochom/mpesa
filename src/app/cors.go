package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/ochom/mpesa/src/app/config"
)

func b2cCors() fiber.Handler {
	crs := cors.ConfigDefault
	crs.AllowOrigins = config.B2CAllowedOrigins

	return cors.New(crs)
}

func taxCors() fiber.Handler {
	crs := cors.ConfigDefault
	crs.AllowOrigins = config.TaxAllowedOrigins
	return cors.New(crs)
}
