package repository

import (
	session "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/repository/session"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const ErrUserNotExists = "User doesn`t exist"

type AuthStorage struct {
	manager *session.SessionManager
	pool    *pgxpool.Pool
}

func NewAuthStorage(pool *pgxpool.Pool, client *redis.Client) *AuthStorage {
	return &AuthStorage{
		manager: session.NewSessionManager(client),
		pool:    pool,
	}
}
