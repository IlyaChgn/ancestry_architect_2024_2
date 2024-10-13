package repository

import (
	"context"
	"fmt"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
	"log"
)

func (storage *TreeStorage) CreateTree(ctx context.Context, userID uint,
	name string) (*models.TreeResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var tree models.TreeResponse

	line := storage.pool.QueryRow(ctx, CreateTreeQuery, userID, name)
	if err := line.Scan(&tree.ID, &tree.UserID, &tree.Name); err != nil {
		customErr := fmt.Errorf("something went wrong while creating tree: %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, customErr
	}

	return &tree, nil
}

func (storage *TreeStorage) AddPermission(ctx context.Context, userID, treeID uint) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	_, err := storage.pool.Exec(ctx, AddPermissionQuery, userID, treeID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while adding permission: %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return customErr
	}

	return nil
}
