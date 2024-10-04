package routers

import (
	nodedel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/node/delivery"
	"github.com/gorilla/mux"
)

func ServeNodeRouter(router *mux.Router, nodeHandler *nodedel.NodeHandler,
	loginRequiredMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/node").Subrouter()
	subrouter.Use(loginRequiredMiddleware)

	subrouter.HandleFunc("/{id:[0-9]+}/description", nodeHandler.GetDescription).Methods("GET")

	subrouter.HandleFunc("/create", nodeHandler.CreateNode).Methods("POST") // tree_id in query
	subrouter.HandleFunc("/{id:[0-9]+}", nodeHandler.DeleteNode).Methods("DELETE")
	subrouter.HandleFunc("/{id:[0-9]+}", nodeHandler.EditNode).Methods("POST")
}
