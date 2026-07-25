package main

import (
	"context"
	"fmt"
	"log"

	api "github.com/cchristian77/wallet-service/entrypoint"
	"github.com/cchristian77/wallet-service/shared/fhttp"
	"github.com/cchristian77/wallet-service/util/config"
	"github.com/cchristian77/wallet-service/util/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zapLog := logger.Initialise()
	defer zapLog.Sync()

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("failed on loading config: %v", err)
	}

	server, err := fhttp.NewHTTPServer(api.InitRouter())
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("failed to create http server: %v", err))
	}

	logger.L().Info(fmt.Sprintf("Starting HTTP Server on Port %d ...", config.Env().App.Port))
	if err = server.Start(ctx); err != nil {
		logger.L().Fatal(fmt.Sprintf("failed to start http server: %v", err))
	}
}
