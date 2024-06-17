package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/mpesa/src/app/config"
	"github.com/ochom/mpesa/src/handlers"
	"github.com/ochom/mpesa/src/models"
)

func init() {
	// init database
	cfg := sql.Config{
		Driver: sql.Driver(config.DbDriver),
		Url:    config.DbUrl,
	}

	if err := sql.New(&cfg); err != nil {
		logs.Fatal("failed to connect to database: %v", err)
	}

	// run migrations
	if err := sql.Conn().AutoMigrate(models.GetSchema()...); err != nil {
		logs.Fatal("failed to run migrations: %v", err)
	}
}

func New() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.ConfigDefault))

	// register routes
	v1 := app.Group("/v1")

	// c2b ...
	c2b := v1.Group("/c2b")
	c2b.Post("/initiate", handlers.HandleStkPush)
	c2b.Post("/result", handlers.HandleStkCallback)
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
	tax.Post("/remit", handlers.HandleTaxRemittance)
	tax.Post("/result", handlers.HandleTaxRemittance)
	tax.Post("/timeout", handlers.HandleTaxRemittance)

	return app
}
