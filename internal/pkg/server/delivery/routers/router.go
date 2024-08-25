package routers

import (
	authdel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/delivery"
	authusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	mylogger "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/log"
	myrecovery "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/recover"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.SugaredLogger,
	authStorage authusecases.AuthStorageInterface) *mux.Router {
	router := mux.NewRouter()
	logMiddleware := mylogger.CreateLogMiddleware(logger)

	router.Use(myrecovery.RecoveryMiddleware)
	router.Use(logMiddleware)

	authHandler := authdel.NewAuthHandler(authStorage)

	rootRouter := router.PathPrefix("/api").Subrouter()
	ServeAuthRouter(rootRouter, authHandler)

	return router
}
