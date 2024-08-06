package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/config"
	myrouter "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery/routers"
	"github.com/gorilla/handlers"
)

type Server struct {
	server *http.Server
}

type serverConfig struct {
	Address string
	Timeout time.Duration
	Handler *http.Handler
}

func createServerConfig(cfg *config.Config, handler *http.Handler) serverConfig {
	return serverConfig{
		Address: cfg.Server.Host + cfg.Server.Port,
		Timeout: time.Second * time.Duration(cfg.Server.Timeout),
		Handler: handler,
	}
}

func createServer(config serverConfig) *http.Server {
	return &http.Server{
		Addr:         config.Address,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		Handler:      *config.Handler,
	}
}

func (srv *Server) Run() error {
	cfgPath := os.Getenv("CONFIG_PATH")
	cfg := config.ReadConfig(cfgPath)

	router := myrouter.NewRouter()

	credentials := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders(cfg.Server.Headers)
	originsOk := handlers.AllowedOrigins(cfg.Server.Origins)
	methodsOk := handlers.AllowedMethods(cfg.Server.Methods)

	muxWithCORS := handlers.CORS(credentials, originsOk, headersOk, methodsOk)(router)

	serverCfg := createServerConfig(cfg, &muxWithCORS)
	srv.server = createServer(serverCfg)

	log.Println("Server is running on port", cfg.Server.Host+cfg.Server.Port)

	return srv.server.ListenAndServe()
}
