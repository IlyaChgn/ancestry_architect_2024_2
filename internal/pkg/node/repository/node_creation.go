package repository

import (
	"context"
	"fmt"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
	"log"
)

func (storage *NodeStorage) CreateNode(ctx context.Context, nodeData *models.CreateNodeRequest) (*models.Node, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var (
		node *models.Node
		err  error
	)

	if nodeData.IsFirstNode {
		node, err = storage.createRootNode(ctx, nodeData)
	} else {
		node, err = storage.createNonRootNode(ctx, nodeData)
	}

	if err != nil {
		customErr := fmt.Errorf("node creation failed: %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	if !nodeData.IsFirstNode {
		err := storage.writeRelatives(ctx, &nodeData.Relatives, node.ID)
		if err != nil {
			customErr := fmt.Errorf("writing relatives failed: %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}

		node.Children = nodeData.Relatives.Children
		node.Spouses = nodeData.Relatives.Spouses
	}

	additional, err := storage.writeAdditionalData(ctx, &nodeData.Addition, node.ID)
	if err != nil {
		customErr := fmt.Errorf("writing additional data failed: %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	node.Birthdate = additional.Birthdate
	node.Deathdate = additional.Deathdate

	return node, nil
}

func (storage *NodeStorage) createRootNode(ctx context.Context,
	nodeData *models.CreateNodeRequest) (*models.Node, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	layerID, err := storage.createLayerIfNotExists(ctx, nodeData.TreeID, 0)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while creating or getting layer, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	node, err := storage.createNode(ctx, layerID, nodeData.Name)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while creating node, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return node, nil
}

func (storage *NodeStorage) createNonRootNode(ctx context.Context,
	nodeData *models.CreateNodeRequest) (*models.Node, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var (
		relative *models.Relative
		err      error
		layerID  uint
	)

	switch {
	case nodeData.Relatives.Children != nil:
		relative, err = storage.getRelativeNode(ctx, nodeData.Relatives.Children[0])
	case nodeData.Relatives.Parents != nil:
		relative, err = storage.getRelativeNode(ctx, nodeData.Relatives.Parents[0])
	case nodeData.Relatives.Spouses != nil:
		relative, err = storage.getRelativeNode(ctx, nodeData.Relatives.Spouses[0])
	default:
		customErr := fmt.Errorf("wrong node data: expected one of siblings, children, or spouse")

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, customErr
	}

	if err != nil {
		customErr := fmt.Errorf("something went wrong while getting relative node, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	if nodeData.Relatives.Children != nil {
		layerID, err = storage.createLayerIfNotExists(ctx, nodeData.TreeID, relative.LayerNumber+1)
	} else if nodeData.Relatives.Parents != nil {
		layerID, err = storage.createLayerIfNotExists(ctx, nodeData.TreeID, relative.LayerNumber-1)
	} else {
		layerID = relative.LayerID
	}

	if err != nil {
		customErr := fmt.Errorf("something went wrong while getting or creating layer, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	node, err := storage.createNode(ctx, layerID, nodeData.Name)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while creating node, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return node, nil
}

func (storage *NodeStorage) createNode(ctx context.Context, layerID uint, name string) (*models.Node, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var node models.Node

	line := storage.pool.QueryRow(ctx, CreateNodeQuery, layerID, name)
	if err := line.Scan(&node.ID, &node.LayerID, &node.Name, &node.Birthdate, &node.Deathdate); err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return &node, nil
}
