package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cchristian77/wallet-service/entrypoint/transfer"
	"github.com/cchristian77/wallet-service/repository"
	"github.com/cchristian77/wallet-service/shared/external/database"
	"github.com/cchristian77/wallet-service/shared/fhttp/middleware"
	"github.com/cchristian77/wallet-service/util/config"
	"github.com/cchristian77/wallet-service/util/logger"
)

// Controller - Interface for controllers to implement for route registration
type Controller interface {
	RegisterRoutes(mux *http.ServeMux)
}

func StartControllers(mux *http.ServeMux) error {
	logger.L().Info("Registering routes for controllers ...")

	ctx := context.Background()

	// initialize DB
	db := database.ConnectToDB()
	if db == nil {
		logger.L().Fatal("Can't connect to Postgres!")
	}

	gormDB, err := database.OpenGormDB(db)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("gorm driver error: %v", err))
	}

	// Initialize repository layer
	repo := repository.NewRepository(gormDB)

	// Initialize controller layer
	transferController, err := transfer.NewController(ctx, repo, gormDB)
	if err != nil {
		return err
	}

	// register routes
	transferController.RegisterRoutes(mux)

	return nil
}

func Initialize() http.Handler {
	router := http.NewServeMux()

	var handler http.Handler = router
	handler = middleware.LogResponse(handler)
	handler = middleware.LogRequest(handler)
	handler = middleware.PanicRecovery(handler)
	handler = middleware.ResponseTime(handler)
	handler = middleware.CorrelationID(handler)

	subRouter := http.NewServeMux()
	if err := StartControllers(subRouter); err != nil {
		logger.L().Fatal(fmt.Sprintf("failed to start controllers: %v", err))
	}

	if apiPrefix := config.Env().App.APIPrefix; apiPrefix != "" {
		router.Handle(apiPrefix+"/", http.StripPrefix(apiPrefix, subRouter))
	} else {
		router.Handle("/", subRouter)
	}

	return handler
}
