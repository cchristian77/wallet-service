package api

import (
	"fmt"
	"net/http"

	transferEntrypoint "github.com/cchristian77/wallet-service/entrypoint/transfer"
	"github.com/cchristian77/wallet-service/repository"
	transactionLedger "github.com/cchristian77/wallet-service/service/transaction_ledger"
	"github.com/cchristian77/wallet-service/shared/external/database"
	"github.com/cchristian77/wallet-service/shared/fhttp"
	"github.com/cchristian77/wallet-service/shared/fhttp/middleware"
	"github.com/cchristian77/wallet-service/util/logger"
)

func InitRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /healthcheck", fhttp.AppHandler(fhttp.DefaultHealthCheckHandler))

	registerRoutes(mux)

	var handler http.Handler = mux
	handler = middleware.LogResponse(handler)
	handler = middleware.LogRequest(handler)
	handler = middleware.PanicRecovery(handler)
	handler = middleware.ResponseTime(handler)
	handler = middleware.CorrelationID(handler)

	return handler
}

func registerRoutes(mux *http.ServeMux) {
	db := database.ConnectToDB()
	if db == nil {
		logger.L().Fatal("Can't connect to Postgres!")
	}

	gormDB, err := database.OpenGormDB(db)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("gorm driver error: %v", err))
	}

	repo := repository.NewRepository(gormDB)

	transactionLedgerService, err := transactionLedger.NewService(repo, gormDB)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("transaction_ledger service initialization error: %v", err))
	}

	transferController := transferEntrypoint.NewController(transactionLedgerService)
	transferController.RegisterRoutes(mux)
}
