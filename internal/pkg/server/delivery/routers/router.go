package routers

import (
	myrecovery "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/recover"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(myrecovery.RecoveryMiddleware)

	// rootRouter := router.PathPrefix("/api").Subrouter()

	return router
}
