package routers

import (
	delivery "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/rest"
	"github.com/gorilla/mux"
)

func ServeAdminAuthRouter(router *mux.Router, adminHandler *delivery.AdminHandler,
	adminRequiredMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/auth").Subrouter()

	subrouter.HandleFunc("/login", adminHandler.Login).Methods("POST")
	subrouter.HandleFunc("/check_auth", adminHandler.CheckAuth).Methods("GET")

	subrouterAdminRequired := subrouter.PathPrefix("").Subrouter()
	subrouterAdminRequired.Use(adminRequiredMiddleware)

	subrouterAdminRequired.HandleFunc("/logout", adminHandler.Logout).Methods("POST")
	subrouterAdminRequired.HandleFunc("/password", adminHandler.EditUserPassword).Methods("POST")
	subrouterAdminRequired.HandleFunc("/list", adminHandler.GetUsersList).Methods("GET")
	subrouterAdminRequired.HandleFunc("/user/{id:[0-9]+}", adminHandler.DeleteUser).Methods("DELETE")
	subrouterAdminRequired.HandleFunc("/user", adminHandler.CreateUser).Methods("POST")
}
