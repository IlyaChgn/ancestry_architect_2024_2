package routers

import (
	treedel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/tree/delivery"
	"github.com/gorilla/mux"
)

func ServeTreeRouter(router *mux.Router, treeHandler *treedel.TreeHandler,
	loginRequiredMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/tree").Subrouter()
	subrouter.Use(loginRequiredMiddleware)

	subrouter.HandleFunc("/{id:[0-9]+}", treeHandler.GetTree).Methods("GET")
	subrouter.HandleFunc("/list/available", treeHandler.GetAvailableTreesList).Methods("GET")
	subrouter.HandleFunc("/list/created", treeHandler.GetCreatedTreesList).Methods("GET")

	subrouter.HandleFunc("", treeHandler.CreateTree).Methods("POST")

	subrouter.HandleFunc("/permission", treeHandler.AddPermission).Methods("POST")
}
