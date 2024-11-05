package routers

import (
	delivery "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/rest"
	"github.com/gorilla/mux"
)

func ServeAdminAuthRouter(router *mux.Router, adminHandler *delivery.AdminHandler) {
	subrouter := router.PathPrefix("/auth").Subrouter()

	subrouter.HandleFunc("/login", adminHandler.Login).Methods("POST")
	subrouter.HandleFunc("/logout", adminHandler.Logout).Methods("POST")
	subrouter.HandleFunc("/password", adminHandler.EditUserPassword).Methods("POST")
	subrouter.HandleFunc("/list", adminHandler.GetUsersList).Methods("GET")
}
