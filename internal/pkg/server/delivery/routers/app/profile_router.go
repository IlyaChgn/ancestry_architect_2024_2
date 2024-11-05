package routers

import (
	profiledel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/delivery"
	"github.com/gorilla/mux"
)

func ServeAppProfileRouter(router *mux.Router, profileHandler *profiledel.ProfileHandler,
	loginRequiredMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/profile").Subrouter()

	subrouterLoginRequired := subrouter.PathPrefix("").Subrouter()
	subrouterLoginRequired.Use(loginRequiredMiddleware)
	subrouterLoginRequired.HandleFunc("", profileHandler.UpdateProfile).Methods("POST")

	subrouter.HandleFunc("/{id:[0-9]+}", profileHandler.GetProfile).Methods("GET")
}
