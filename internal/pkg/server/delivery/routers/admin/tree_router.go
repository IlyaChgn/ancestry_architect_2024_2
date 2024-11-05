package routers

import (
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/rest"
	"github.com/gorilla/mux"
)

func ServeAdminTreeRouter(router *mux.Router, adminHandler *delivery.AdminHandler) {
	subrouter := router.PathPrefix("/tree").Subrouter()

	subrouter.HandleFunc("/list", adminHandler.GetTreesList).Methods("GET")
	subrouter.HandleFunc("/name", adminHandler.EditTreeName).Methods("POST")
	subrouter.HandleFunc("/{id:[0-9]+}/nodes", adminHandler.GetNodesList).Methods("GET")
}
