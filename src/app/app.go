package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ochom/mpesa/src/handlers"
)

func New() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.ConfigDefault))

	// serve swagger docs at root
	app.Static("/", "./docs")

	// register routes
	v1 := app.Group("/v1")
	v1.Route("/accounts", func(r fiber.Router) {
		r.Use(basicAuth())
		r.Get("/", handlers.HandleListAccounts)
		r.Get("/search", handlers.HandleSearchAccounts)
		r.Post("/", handlers.HandleCreateAccount)
		r.Put("/:id", handlers.HandleUpdateAccount)
		r.Delete("/:id", handlers.HandleDeleteAccount)
		r.Post("/register-urls", handlers.HandleC2BRegisterUrls)
	})

	// c2b ...
	v1.Route("/c2b", func(r fiber.Router) {
		r.Get("/payments", handlers.HandleGetC2BPayments)
		r.Post("/initiate", handlers.HandleStkPush)
		r.Post("/result", safOrigins(), handlers.HandleC2BCallback)
		r.Post("/validate", safOrigins(), handlers.HandleRestValidation)
		r.Post("/confirm", safOrigins(), handlers.HandleRestConfirmation)
		r.Post("/soap/validate", handlers.HandleSoapValidation)
		r.Post("/soap/confirm", handlers.HandleSoapConfirmation)
	})

	// b2c ...
	v1.Route("/b2c", func(r fiber.Router) {
		r.Get("/payments", handlers.HandleGetB2CPayments)
		r.Post("/initiate", b2cOrigins(), handlers.HandleInitiatePayment)
		r.Post("/result", safOrigins(), handlers.HandleB2CResult)
		r.Post("/timeout", safOrigins(), handlers.HandleB2CTimeout)
	})

	// tax ...
	v1.Route("/tax", func(r fiber.Router) {
		r.Post("/initiate", taxOrigins(), handlers.HandleTaxRemittance)
		r.Post("/result", safOrigins(), handlers.HandleTaxResult)
		r.Post("/timeout", safOrigins(), handlers.HandleTaxTimeout)
	})

	return app
}
