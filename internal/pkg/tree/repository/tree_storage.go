package repository

import pool "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/repository"

type TreeStorage struct {
	pool pool.PostgresPool
}

func NewTreeStorage(pool pool.PostgresPool) *TreeStorage {
	return &TreeStorage{
		pool: pool,
	}
}
