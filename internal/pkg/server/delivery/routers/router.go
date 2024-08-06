package routers

import (
	mylogger "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/log"
	myrecovery "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/recover"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.SugaredLogger) *mux.Router {
	router := mux.NewRouter()
	logMiddleware := mylogger.CreateLogMiddleware(logger)

	router.Use(myrecovery.RecoveryMiddleware)
	router.Use(logMiddleware)

	// rootRouter := router.PathPrefix("/api").Subrouter()

	return router
}
