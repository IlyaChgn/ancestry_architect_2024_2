package repository

import (
	"context"
	"fmt"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
	"log"
	"time"
)

func (storage *AdminStorage) GetTreesList(ctx context.Context) (*[]models.TreeResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var list []models.TreeResponse

	rows, err := storage.pool.Query(ctx, GetTreesListQuery)
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
			customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, customErr
		}

		list = append(list, tree)
	}

	return &list, nil
}

func (storage *AdminStorage) GetTreesListByUserID(ctx context.Context, userID uint) (*[]models.TreeResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var list []models.TreeResponse

	rows, err := storage.pool.Query(ctx, GetTreesListByUserIDQuery, userID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while getting trees list by user ID, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, customErr
	}

	defer rows.Close()

	for rows.Next() {
		var tree models.TreeResponse

		if err := rows.Scan(&tree.ID, &tree.UserID, &tree.Name); err != nil {
			customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, customErr
		}

		list = append(list, tree)
	}

	return &list, nil
}

func (storage *AdminStorage) GetNodesList(ctx context.Context, treeID uint) (*[]models.NodeForAdmin, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var list []models.NodeForAdmin

	rows, err := storage.pool.Query(ctx, GetNodesListQuery, treeID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while getting nodes list, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, customErr
	}

	defer rows.Close()

	for rows.Next() {
		var (
			node        models.NodeForAdmin
			birthdate   *time.Time
			deathdate   *time.Time
			previewPath *string
		)

		if err := rows.Scan(&node.ID, &node.Name, &birthdate, &deathdate, &node.Gender, &previewPath,
			&node.IsDeleted, &node.LayerID, &node.LayerNum, &node.TreeID, &node.UserID); err != nil {
			customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, customErr
		}

		if birthdate != nil {
			node.Birthdate = *birthdate
		}

		if deathdate != nil {
			node.Deathdate = *deathdate
		}

		if previewPath != nil {
			node.PreviewPath = *previewPath
		}

		list = append(list, node)
	}

	return &list, nil
}

func (storage *AdminStorage) EditTreeName(ctx context.Context, treeID uint,
	name string) (*models.TreeResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var tree models.TreeResponse

	line := storage.pool.QueryRow(ctx, EditTreeNameQuery, treeID, name)
	if err := line.Scan(&tree.ID, &tree.UserID, &tree.Name); err != nil {
		customErr := fmt.Errorf("something went wrong while editing tree, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, customErr
	}

	return &tree, nil
}
