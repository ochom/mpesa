package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/ochom/mpesa/src/handlers"
)

func New() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.ConfigDefault))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello Broker")
	})

	// register routes
	v1 := app.Group("/v1")

	// c2b ...
	c2b := v1.Group("/c2b")
	c2b.Post("/register-urls", handlers.HandleC2BRegisterUrls)
	c2b.Post("/initiate", handlers.HandleStkPush)
	c2b.Post("/result", handlers.HandleC2BResult)
	c2b.Post("/rest/validate", handlers.HandleRestValidation)
	c2b.Post("/rest/confirm", handlers.HandleRestConfirmation)
	c2b.Post("/soap/validate", handlers.HandleSoapValidation)
	c2b.Post("/soap/confirm", handlers.HandleSoapConfirmation)

	// b2c ...
	b2c := v1.Group("b2c")
	b2c.Use(b2cCors())
	b2c.Post("/initiate", handlers.HandleInitiatePayment)
	b2c.Post("/result", handlers.HandleB2CResult)
	b2c.Post("/timeout", handlers.HandleB2cTimeout)

	// tax ...
	tax := v1.Group("/tax")
	tax.Use(taxCors())
	tax.Post("/initiate", handlers.HandleTaxRemittance)
	tax.Post("/result", handlers.HandleTaxResult)
	tax.Post("/timeout", handlers.HandleTaxTimeout)

	return app
}
