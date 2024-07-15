package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/mpesa/src/app"
	"github.com/ochom/mpesa/src/models"
)

// init logger
func init() {
	logPath := env.Get("LOGS_PATH", "./data/logs.txt")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logs.Error("[❌] Failed to open log file: %s", err.Error())
		return
	}

	logWriter := io.MultiWriter(os.Stdout, f)
	logs.SetOutput(logWriter)
}

func init() {
	// init database
	cfg := sql.Config{
		Driver: env.Int("DATABASE_DRIVER", 0),
		Url:    env.Get("DATABASE_URL"),
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
