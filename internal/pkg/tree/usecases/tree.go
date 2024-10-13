package usecases

import (
	"context"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
)

type TreeStorageInterface interface {
	CheckPermission(ctx context.Context, treeID, userID uint) (bool, error)

	GetCreatedTrees(ctx context.Context, userID uint) ([]*models.TreeResponse, error)
	GetAvailableTrees(ctx context.Context, userID uint) ([]*models.TreeResponse, error)

	CreateTree(ctx context.Context, userID uint, name string) (*models.TreeResponse, error)
	AddPermission(ctx context.Context, userID, treeID uint) error
}
