package repository

import (
	pool "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/repository"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/session"
	"github.com/redis/go-redis/v9"
)

type AdminStorage struct {
	manager *session.SessionManager
	pool    pool.PostgresPool
}

func NewAdminStorage(pool pool.PostgresPool, client *redis.Client) *AdminStorage {
	return &AdminStorage{
		manager: session.NewSessionManager(client),
		pool:    pool,
	}
}
