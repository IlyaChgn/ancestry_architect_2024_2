package routers

import (
	nodedel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/node/delivery"
	"github.com/gorilla/mux"
)

func ServeNodeRouter(router *mux.Router, nodeHandler *nodedel.NodeHandler,
	loginRequiredMiddleware, permissionRequiredMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/node").Subrouter()
	subrouter.Use(loginRequiredMiddleware)

	subrouterPermissionRequired := subrouter.PathPrefix("/{id:[0-9]+}").Subrouter()
	subrouterPermissionRequired.Use(permissionRequiredMiddleware)

	subrouterPermissionRequired.HandleFunc("/description", nodeHandler.GetDescription).Methods("GET")
	subrouterPermissionRequired.HandleFunc("", nodeHandler.DeleteNode).Methods("DELETE")
	subrouterPermissionRequired.HandleFunc("", nodeHandler.EditNode).Methods("POST")
	subrouterPermissionRequired.HandleFunc("/preview", nodeHandler.UpdateAvatar).Methods("POST")

	subrouter.HandleFunc("", nodeHandler.CreateNode).Methods("POST")
}
