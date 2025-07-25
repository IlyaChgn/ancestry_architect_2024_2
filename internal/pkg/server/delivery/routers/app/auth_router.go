package routers

import (
	authdel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/delivery"
	"github.com/gorilla/mux"
)

func ServeAppAuthRouter(router *mux.Router, authHandler *authdel.AuthHandler,
	loginRequiredMiddleware, logoutRequiredMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/auth").Subrouter()

	subrouterLoginRequired := subrouter.PathPrefix("").Subrouter()
	subrouterLoginRequired.Use(loginRequiredMiddleware)
	subrouterLoginRequired.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	subrouterLogoutRequired := subrouter.PathPrefix("").Subrouter()
	subrouterLogoutRequired.Use(logoutRequiredMiddleware)
	subrouterLogoutRequired.HandleFunc("/login", authHandler.Login).Methods("POST")
	subrouterLogoutRequired.HandleFunc("/signup", authHandler.Signup).Methods("POST")

	subrouter.HandleFunc("/check_auth", authHandler.CheckAuth).Methods("GET")
}
