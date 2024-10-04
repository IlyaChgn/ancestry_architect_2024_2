package server

import (
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/node/repository"
	treerepo "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/tree/repository"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	authrepo "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/repository"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/config"
	profilerepo "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/repository"
	myrouter "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery/routers"
	pool "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/repository"
	logging "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/usecases"
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

func createServerConfig(addr string, timeout int, handler *http.Handler) serverConfig {
	return serverConfig{
		Address: addr,
		Timeout: time.Second * time.Duration(timeout),
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
	if cfg == nil {
		log.Fatal("The config wasn`t opened")
	}

	logger, err := logging.NewLogger(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "))
	if err != nil {
		log.Fatal("Something went wrong while creating logger: ", err)
	}

	postgresURL := pool.NewConnectionString(cfg.Postgres.Username, cfg.Postgres.Password,
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)

	postgresPool, err := pool.NewPostgresPool(postgresURL)
	if err != nil {
		log.Fatal("Something went wrong while connecting to postgres database: ", err)
	}

	redisClient := pool.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)

	authStorage := authrepo.NewAuthStorage(postgresPool, redisClient)
	profileStorage := profilerepo.NewProfileStorage(postgresPool, os.Getenv("STATIC_DIRECTORY"))
	treeStorage := treerepo.NewTreeStorage(postgresPool)
	nodeStorage := repository.NewNodeStorage(postgresPool, os.Getenv("STATIC_DIRECTORY"))

	router := myrouter.NewRouter(logger, authStorage, profileStorage, treeStorage, nodeStorage)

	credentials := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders(cfg.Server.Headers)
	originsOk := handlers.AllowedOrigins(cfg.Server.Origins)
	methodsOk := handlers.AllowedMethods(cfg.Server.Methods)

	muxWithCORS := handlers.CORS(credentials, originsOk, headersOk, methodsOk)(router)

	serverCfg := createServerConfig(cfg.Server.Host+cfg.Server.Port, cfg.Server.Timeout, &muxWithCORS)
	srv.server = createServer(serverCfg)

	log.Printf("Server is listening on %s\n", cfg.Server.Host+cfg.Server.Port)

	return srv.server.ListenAndServe()
}
