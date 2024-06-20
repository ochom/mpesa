package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/mpesa/src/app"
	"github.com/ochom/mpesa/src/app/config"
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

func main() {
	srv := app.New()
	go func() {
		port := helpers.GetPort(8080)
		if err := srv.Listen(fmt.Sprintf(":%d", port)); err != nil {
			logs.Fatal("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.ShutdownWithContext(ctx); err != nil {
		logs.Fatal("failed to shutdown server: %v", err)
	}
}
