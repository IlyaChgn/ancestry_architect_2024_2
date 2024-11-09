package usecases

import (
	"context"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
)

type AdminStorageInterface interface {
	GetAdminByEmail(ctx context.Context, email string) (*models.User, error)
	GetAdminBySessionID(ctx context.Context, sessionID string) (*models.User, error)
	UpdatePassword(ctx context.Context, id uint, password string) (*models.User, error)
	GetUsersList(ctx context.Context) (*[]models.User, error)
	CreateSession(ctx context.Context, sessionID string, userID uint) error
	RemoveSession(ctx context.Context, sessionID string) error

	GetTreesList(ctx context.Context) (*[]models.TreeResponse, error)
	GetTreesListByUserID(ctx context.Context, userID uint) (*[]models.TreeResponse, error)
	GetNodesList(ctx context.Context, treeID uint) (*[]models.NodeForAdmin, error)
	EditTreeName(ctx context.Context, treeID uint, name string) (*models.TreeResponse, error)
}
