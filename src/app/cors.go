package app

import (
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/mpesa/src/app/config"
)

func getCors() cors.Config {
	crs := cors.ConfigDefault
	crs.AllowOrigins = ""
	return crs
}

func basicAuth() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			config.BasicAuthUsername: config.BasicAuthPassword,
		},
		Unauthorized: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		},
	})
}

func safOrigins() fiber.Handler {
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

	crs := getCors()
	crs.AllowOriginsFunc = func(origin string) bool {
		logs.Info("Receive request from origin: %s", origin)
		return slices.Contains(allowedOrigins, origin)
	}

	return cors.New(crs)
}

func b2cOrigins() fiber.Handler {
	allowedOrigins := config.B2CAllowedOrigins

	crs := getCors()
	crs.AllowOriginsFunc = func(origin string) bool {
		logs.Info("Receive request from origin: %s", origin)
		if allowedOrigins == "" || allowedOrigins == "*" {
			return true
		}

		return strings.Contains(allowedOrigins, origin)
	}
	return cors.New(crs)
}

func taxOrigins() fiber.Handler {
	allowedOrigins := config.TaxAllowedOrigins

	crs := getCors()
	crs.AllowOriginsFunc = func(origin string) bool {
		logs.Info("Receive request from origin: %s", origin)
		if allowedOrigins == "" || allowedOrigins == "*" {
			return true
		}

		return strings.Contains(allowedOrigins, origin)
	}
	return cors.New(crs)
}
