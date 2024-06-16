package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/mpesa/src/handlers"
	"github.com/ochom/mpesa/src/models"
)

func init() {
	// init database
	cfg := sql.Config{
		DatabaseType: sql.MySQL,
		Url:          env.Get("DATABASE_URL"),
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
	c2b := v1.Group("c2b")
	c2b.Post("initiate", handlers.HandleStkPush)
	c2b.Post("callback", handlers.HandleC2bConfirmation)
	c2b.Post("validate", handlers.HandleC2bValidation)
	c2b.Post("confirm", handlers.HandleC2bConfirmation)

	// b2c ...
	b2c := v1.Group("b2c")
	b2c.Use(b2cCors())
	b2c.Post("initiate", handlers.HandleInitiatePayment)
	b2c.Post("callback", handlers.HandlePaymentCallback)

	return app
}
