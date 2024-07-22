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
	sc := v1.Group("/accounts")
	sc.Use(basicAuth())
	sc.Get("/", handlers.HandleListAccounts)
	sc.Get("/search", handlers.HandleSearchAccounts)
	sc.Post("/", handlers.HandleCreateAccount)
	sc.Put("/:id", handlers.HandleUpdateAccount)
	sc.Delete("/:id", handlers.HandleDeleteAccount)
	sc.Post("/register-urls", handlers.HandleC2BRegisterUrls)

	// c2b ...
	c2b := v1.Group("/c2b")
	c2b.Get("/payments", handlers.HandleGetC2BPayments)
	c2b.Post("/initiate", handlers.HandleStkPush)
	c2b.Post("/result", safOrigins(), handlers.HandleC2BResult)
	c2b.Post("/validate", safOrigins(), handlers.HandleRestValidation)
	c2b.Post("/confirm", safOrigins(), handlers.HandleRestConfirmation)
	c2b.Post("/soap/validate", handlers.HandleSoapValidation)
	c2b.Post("/soap/confirm", handlers.HandleSoapConfirmation)

	// b2c ...
	b2c := v1.Group("b2c")
	b2c.Get("/payments", handlers.HandleGetB2CPayments)
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
