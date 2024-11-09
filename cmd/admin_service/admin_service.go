package main

import (
	"fmt"
	mygrpc "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc"
	adminproto "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/repository"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/config"
	pool "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/repository"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env file", err)
	}

	cfg := &config.AdminConfig{}

	cfgPath := os.Getenv("CONFIG_PATH")
	config.ReadConfig(cfgPath, cfg)
	if cfg == nil {
		log.Fatal("The config wasn`t opened")
	}

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error occurred while listening admin service. %v", err)
	}

	grpcConn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Error occurred while starting grpc connection on admin service, %v", err)
	}
	defer grpcConn.Close()

	postgresURL := pool.NewConnectionString(cfg.Postgres.Username, cfg.Postgres.Password,
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)

	postgresPool, err := pool.NewPostgresPool(postgresURL)
	if err != nil {
		log.Fatal("Something went wrong while connecting to postgres database: ", err)
	}

	redisClient := pool.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)

	adminStorage := repository.NewAdminStorage(postgresPool, redisClient)
	adminManager := mygrpc.NewAdminManager(adminStorage)

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor()))
	adminproto.RegisterAdminServer(srv, adminManager)

	log.Printf("Admin service is listening on %s\n", addr)

	err = srv.Serve(listener)
	if err != nil {
		log.Fatalf("error occurred while serving grpc. %v", err)
	}
}
