package routers

import (
	authdel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/delivery"
	"github.com/gorilla/mux"
)

func ServeAuthRouter(router *mux.Router, authHandler *authdel.AuthHandler) {
	subrouter := router.PathPrefix("/auth").Subrouter()

	subrouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	subrouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	subrouter.HandleFunc("/check_auth", authHandler.CheckAuth).Methods("GET")
	subrouter.HandleFunc("/signup", authHandler.Signup).Methods("POST")
}
