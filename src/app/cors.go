package app

import (
	"path"
	"slices"
	"strings"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/logs"
)

func docs() func(*fiber.Ctx) error {
	rootPath := env.Get("ROOT_PATH", "/")
	return swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.yaml",
		Path:     path.Join(rootPath, "docs"),
		Title:    "Mpesa Broker API",
	})
}

func getCors() cors.Config {
	crs := cors.ConfigDefault
	crs.AllowOrigins = ""
	return crs
}

func basicAuth() fiber.Handler {
	username := env.Get("BASIC_AUTH_USERNAME")
	password := env.Get("BASIC_AUTH_PASSWORD")
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			username: password,
		},
		Unauthorized: func(c *fiber.Ctx) error {
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
	crs := getCors()

	crs.AllowOriginsFunc = func(origin string) bool {
		logs.Info("Receive request from origin: %s", origin)

		allowedOrigins := env.Get("B2C_ALLOWED_ORIGINS")
		if allowedOrigins == "" || allowedOrigins == "*" {
			return true
		}

		return strings.Contains(allowedOrigins, origin)
	}
	return cors.New(crs)
}

func taxOrigins() fiber.Handler {
	crs := getCors()

	crs.AllowOriginsFunc = func(origin string) bool {
		logs.Info("Receive request from origin: %s", origin)

		allowedOrigins := env.Get("TAX_REMITTANCE_ALLOWED_ORIGINS")
		if allowedOrigins == "" || allowedOrigins == "*" {
			return true
		}

		return strings.Contains(allowedOrigins, origin)
	}
	return cors.New(crs)
}
