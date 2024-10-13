package usecases

import (
	"context"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"mime/multipart"
)

type NodeStorageInterface interface {
	CheckPermission(ctx context.Context, nodeID, userID uint) (bool, error)

	GetDescription(ctx context.Context, nodeID uint) (*models.DescriptionResponse, error)

	CreateNode(ctx context.Context, node *models.CreateNodeRequest) (*models.Node, error)
	DeleteNode(ctx context.Context, nodeID uint) error
	UpdatePreview(ctx context.Context, preview *multipart.FileHeader, nodeID uint) (string, error)
}
