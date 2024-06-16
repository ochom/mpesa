package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/ochom/gutils/logs"
	"github.com/ochom/mpesa/src/app"
)

func main() {
	srv := app.New()
	go func() {
		if err := srv.Listen(":8080"); err != nil {
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
