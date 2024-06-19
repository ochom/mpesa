package app

import (
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/ochom/mpesa/src/app/config"
)

func safOrigins(next fiber.Handler) fiber.Handler {
	origins := []string{
		"196.201.214.200",
		"196.201.214.206",
		"196.201.214.207",
		"196.201.214.208",
		"196.201.213.114",
		"196.201.213.44",
		"196.201.212.127",
		"196.201.212.138",
		"196.201.212.129",
		"196.201.212.136",
		"196.201.212.74",
		"196.201.212.69",
	}

	return func(c fiber.Ctx) error {
		origin := c.Get("Origin")
		if slices.Contains(origins, origin) {
			return next(c)
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
}

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
