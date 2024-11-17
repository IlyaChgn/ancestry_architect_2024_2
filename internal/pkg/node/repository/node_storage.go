package repository

import pool "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/repository"

const (
	Parent = "Родитель"
	Spouse = "Супруг"
)

type NodeStorage struct {
	pool            pool.PostgresPool
	staticDirectory string
}

func NewNodeStorage(pool pool.PostgresPool, staticDirectory string) *NodeStorage {
	return &NodeStorage{
		pool:            pool,
		staticDirectory: staticDirectory,
	}
}
