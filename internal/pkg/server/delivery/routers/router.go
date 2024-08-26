package routers

import (
	authdel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/delivery"
	authusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	myauth "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/auth"
	mylogger "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/log"
	myrecovery "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/recover"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.SugaredLogger,
	authStorage authusecases.AuthStorageInterface) *mux.Router {
	router := mux.NewRouter()

	router.Use(myrecovery.RecoveryMiddleware)
	router.Use(mylogger.CreateLogMiddleware(logger))

	authHandler := authdel.NewAuthHandler(authStorage)

	loginRequiredMiddleware := myauth.LoginRequiredMiddleware(authStorage)
	logoutRequiredMiddleware := myauth.LogoutRequiredMiddleware(authStorage)

	rootRouter := router.PathPrefix("/api").Subrouter()
	ServeAuthRouter(rootRouter, authHandler, loginRequiredMiddleware, logoutRequiredMiddleware)

	return router
}
