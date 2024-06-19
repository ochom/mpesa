package app

import (
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/mpesa/src/app/config"
)

func safOrigins(next fiber.Handler) fiber.Handler {
	allowedOrigins := []string{
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
		if slices.Contains(allowedOrigins, origin) {
			return next(c)
		}

		logs.Warn("received request from unknown origin: %s", origin)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
}

func b2cOrigins(next fiber.Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		allowedOrigins := config.B2CAllowedOrigins
		if allowedOrigins == "" {
			allowedOrigins = "*"
		}

		if allowedOrigins == "*" {
			return next(c)
		}

		origin := c.Get("Origin")
		if strings.Contains(allowedOrigins, origin) {
			return next(c)
		}

		logs.Warn("received b2c request from unknown origin: %s", origin)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
}

func taxOrigins(next fiber.Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		allowedOrigins := config.TaxAllowedOrigins
		if allowedOrigins == "" {
			allowedOrigins = "*"
		}

		if allowedOrigins == "*" {
			return c.Next()
		}

		origin := c.Get("Origin")
		if strings.Contains(allowedOrigins, origin) {
			return next(c)
		}

		logs.Warn("received tax request from unknown origin: %s", origin)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
}
