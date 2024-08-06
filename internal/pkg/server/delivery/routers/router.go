package routers

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	// rootRouter := router.PathPrefix("/api").Subrouter()

	return router
}
