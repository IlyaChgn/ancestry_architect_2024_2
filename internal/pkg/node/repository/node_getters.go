package repository

import (
	"context"
	"fmt"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
	"log"
)

func (storage *NodeStorage) CheckPermission(ctx context.Context, nodeID, userID uint) (bool, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var hasPermission bool

	line := storage.pool.QueryRow(ctx, CheckPermissionForNodeQuery, nodeID, userID)
	if err := line.Scan(&hasPermission); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return false, err
	}

	return hasPermission, nil
}

func (storage *NodeStorage) GetDescription(ctx context.Context, nodeID uint) (*models.DescriptionResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var description models.DescriptionResponse

	line := storage.pool.QueryRow(ctx, GetDescriptionQuery, nodeID)
	if err := line.Scan(&description.NodeID, &description.Description); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return &description, nil
}

func (storage *NodeStorage) getRelativeNode(ctx context.Context, relativeID uint) (*models.Relative, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var relative models.Relative

	line := storage.pool.QueryRow(ctx, GetRelativeNodeQuery, relativeID)
	if err := line.Scan(&relative.ID, &relative.LayerNumber, &relative.LayerID); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return &relative, nil
}

func (storage *NodeStorage) getParentsList(ctx context.Context, childID uint) ([]uint, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var parentsList []uint

	rows, err := storage.pool.Query(ctx, GetParentsQuery, childID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while getting parents, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	for rows.Next() {
		var parentID uint

		if err := rows.Scan(&parentID); err != nil {
			customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}

		parentsList = append(parentsList, parentID)
	}

	return parentsList, nil
}

func (storage *NodeStorage) getLayer(ctx context.Context, treeID, number uint) (uint, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var layerID uint

	line := storage.pool.QueryRow(ctx, GetLayerQuery, treeID, number)
	if err := line.Scan(&layerID); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return 0, err
	}

	return layerID, nil
}
