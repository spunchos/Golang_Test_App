package main

import (
	"TestApp/internal/app"
	"TestApp/internal/config"
	"TestApp/pkg/logging"
	"context"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("starting application")

	cfg := config.GetConfig()

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	application, err := app.NewApp(ctx, cfg)
	if err != nil {
		logger.Fatalf("failed to initialize application: %v", err)
	}

	application.Run()
}
