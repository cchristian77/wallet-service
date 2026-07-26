package main

import (
	"context"
	"fmt"

	api "github.com/cchristian77/wallet-service/entrypoint"
	"github.com/cchristian77/wallet-service/shared/fhttp"
	"github.com/cchristian77/wallet-service/util/config"
	"github.com/cchristian77/wallet-service/util/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zapLog := logger.Initialize()
	defer zapLog.Sync()

	if err := config.LoadEnv(); err != nil {
		logger.L().Fatal(fmt.Sprintf("failed on loading env : %v", err))
		return
	}

	server, err := fhttp.NewHTTPServer(api.Initialize())
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("failed to create http server : %v", err))
		return
	}

	logger.L().Info(fmt.Sprintf("Starting HTTP Server on Port %d ...", config.Env().App.Port))
	if err = server.Start(ctx); err != nil {
		logger.L().Fatal(fmt.Sprintf("failed to start http server, err: %v", err))
	}
}
