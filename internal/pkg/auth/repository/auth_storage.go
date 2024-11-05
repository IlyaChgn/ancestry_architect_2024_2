package repository

import (
	pool "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/repository"
	session "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/session"
	"github.com/redis/go-redis/v9"
)

const (
	ErrUserNotExists     = "User doesn`t exist"
	ErrUserAlreadyExists = "User with same email already exists"
	ErrWrongEmailFormat  = "Wrong email format"
)

type AuthStorage struct {
	manager *session.SessionManager
	pool    pool.PostgresPool
}

func NewAuthStorage(pool pool.PostgresPool, client *redis.Client) *AuthStorage {
	return &AuthStorage{
		manager: session.NewSessionManager(client),
		pool:    pool,
	}
}
