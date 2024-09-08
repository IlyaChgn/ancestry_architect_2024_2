package repository

import pool "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/repository"

type ProfileStorage struct {
	pool pool.PostgresPool
}

func NewProfileStorage(pool pool.PostgresPool) *ProfileStorage {
	return &ProfileStorage{
		pool: pool,
	}
}
