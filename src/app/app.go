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
	sc := v1.Group("/accounts")
	sc.Use(basicAuth())
	sc.Get("/", handlers.HandleListShortCodes)
	sc.Post("/", handlers.HandleCreateShortCode)
	sc.Put("/:id", handlers.HandleUpdateShortCode)
	sc.Post("/register-urls", handlers.HandleC2BRegisterUrls)

	// c2b ...
	c2b := v1.Group("/c2b")
	c2b.Post("/initiate", handlers.HandleStkPush)
	c2b.Post("/result", safOrigins(), handlers.HandleC2BResult)
	c2b.Post("/validate", safOrigins(), handlers.HandleRestValidation)
	c2b.Post("/confirm", safOrigins(), handlers.HandleRestConfirmation)
	c2b.Post("/soap/validate", handlers.HandleSoapValidation)
	c2b.Post("/soap/confirm", handlers.HandleSoapConfirmation)

	// b2c ...
	b2c := v1.Group("b2c")
	b2c.Post("/initiate", b2cOrigins(), handlers.HandleInitiatePayment)
	b2c.Post("/result", safOrigins(), handlers.HandleB2CResult)
	b2c.Post("/timeout", safOrigins(), handlers.HandleB2cTimeout)

	// tax ...
	tax := v1.Group("/tax")
	tax.Post("/initiate", taxOrigins(), handlers.HandleTaxRemittance)
	tax.Post("/result", safOrigins(), handlers.HandleTaxResult)
	tax.Post("/timeout", safOrigins(), handlers.HandleTaxTimeout)

	return app
}
