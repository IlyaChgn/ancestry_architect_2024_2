package repository

import (
	"context"
	"fmt"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
	"log"
)

func (storage *TreeStorage) CheckPermission(ctx context.Context, treeID, userID uint) (bool, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var hasPermission bool

	line := storage.pool.QueryRow(ctx, CheckPermissionForTreeQuery, treeID, userID)
	if err := line.Scan(&hasPermission); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return false, customErr
	}

	return hasPermission, nil
}

func (storage *TreeStorage) GetCreatedTrees(ctx context.Context, userID uint) ([]*models.TreeResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var treesList []*models.TreeResponse

	rows, err := storage.pool.Query(ctx, GetCreatedTreesListQuery, userID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while getting trees list, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, customErr
	}

	defer rows.Close()

	for rows.Next() {
		var tree models.TreeResponse

		if err := rows.Scan(&tree.ID, &tree.UserID, &tree.Name); err != nil {
			customErr := fmt.Errorf("something went wrong while scanning row, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, customErr
		}

		treesList = append(treesList, &tree)
	}

	return treesList, nil
}

func (storage *TreeStorage) GetAvailableTrees(ctx context.Context, userID uint) ([]*models.TreeResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var treesList []*models.TreeResponse

	rows, err := storage.pool.Query(ctx, GetAvailableTreesListQuery, userID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while getting trees list, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, customErr
	}

	defer rows.Close()

	for rows.Next() {
		var tree models.TreeResponse

		if err := rows.Scan(&tree.ID, &tree.UserID, &tree.Name); err != nil {
			customErr := fmt.Errorf("something went wrong while scanning row, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, customErr
		}

		treesList = append(treesList, &tree)
	}

	return treesList, nil
}

func (storage *TreeStorage) GetTree(ctx context.Context, treeID uint) (*models.Tree, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var tree models.Tree

	rows, err := storage.pool.Query(ctx, GetTreeQuery, treeID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while getting tree, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, customErr
	}

	var (
		treeData   models.Tree
		oldLayerID uint
	)

	defer rows.Close()

	for rows.Next() {
		var (
			layerData   models.Layer
			nodeData    models.Node
			previewPath *string
		)

		if err := rows.Scan(&treeData.ID, &treeData.Name, &layerData.ID, &layerData.Number, &nodeData.ID,
			&nodeData.LayerID, &nodeData.Name, &nodeData.Birthdate, &nodeData.Deathdate, &previewPath,
			&nodeData.IsSpouse, &nodeData.Gender, &nodeData.Relatives); err != nil {
			customErr := fmt.Errorf("something went wrong while scanning row, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, customErr
		}

		if oldLayerID != layerData.ID {
			oldLayerID = layerData.ID

			tree.Layers = append(tree.Layers, layerData)
		}

		if previewPath != nil {
			nodeData.PreviewPath = *previewPath
		}

		tree.Layers[layerData.Number].Nodes = append(tree.Layers[layerData.Number].Nodes, nodeData)
	}

	tree.ID = treeData.ID
	tree.Name = treeData.Name

	return &tree, nil
}
