package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"log"
	"mime/multipart"
	"os"
	"time"
)

func (storage *NodeStorage) UpdatePreview(ctx context.Context, preview *multipart.FileHeader,
	nodeID uint) (string, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var oldPath string

	oldPathLine := storage.pool.QueryRow(ctx, GetPreviewQuery, nodeID)
	if oldPathErr := oldPathLine.Scan(&oldPath); oldPathErr == nil {
		os.Remove(oldPath)
	}

	newPath, err := utils.WriteFile(preview, storage.staticDirectory, "preview")
	if err != nil {
		customErr := fmt.Errorf("something went wrong while writing preview, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return "", err
	}

	_, err = storage.pool.Exec(ctx, UpdatePreviewQuery, newPath, nodeID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return "", err
	}

	return newPath, nil
}

func (storage *NodeStorage) DeleteNode(ctx context.Context, nodeID uint) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	_, err := storage.pool.Exec(ctx, DeleteNodeQuery, nodeID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while seleting node, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return err
	}

	return nil
}

func (storage *NodeStorage) EditNode(ctx context.Context, editedNode *models.EditNodeRequest,
	nodeID uint) (*models.EditNodeResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var (
		err  error
		node models.EditNodeResponse
	)

	if editedNode.Name != "" {
		err = storage.updateName(ctx, editedNode.Name, nodeID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while updating name, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}

		node.Name = editedNode.Name
	}

	if editedNode.Birthdate != "" {
		birthdate, err := storage.updateBirthdate(ctx, editedNode.Birthdate, nodeID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while updating birthdate, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}

		node.Birthdate = birthdate
	}

	if editedNode.Deathdate != "" {
		deathdate, err := storage.updateDeathdate(ctx, editedNode.Deathdate, nodeID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while updating deathdate, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}

		node.Deathdate = deathdate
	}

	if editedNode.Description != "" {
		err = storage.updateDescription(ctx, editedNode.Description, nodeID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while updating description, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}
	}

	if editedNode.Gender != "" {
		err = storage.updateGender(ctx, editedNode.Gender, nodeID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while updating gender, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}

		node.Gender = editedNode.Gender
	}

	return &node, nil
}

func (storage *NodeStorage) writeRelatives(ctx context.Context, relatives *models.GetRelativesList,
	nodeID uint) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	if relatives.Children != nil {
		for _, childID := range relatives.Children {
			err := storage.setRelative(ctx, nodeID, childID, Parent)
			if err != nil {
				customErr := fmt.Errorf("something went wrong while writing children, %v", err)

				utils.LogError(logger, customErr)
				log.Println(customErr)

				return err
			}
		}
	}

	if relatives.Parents != nil {
		for _, parentID := range relatives.Parents {
			err := storage.setRelative(ctx, parentID, nodeID, Parent)
			if err != nil {
				customErr := fmt.Errorf("something went wrong while writing parents, %v", err)

				utils.LogError(logger, customErr)
				log.Println(customErr)

				return err
			}
		}
	}

	if relatives.Spouses != nil {
		for _, spouseID := range relatives.Spouses {
			err := storage.setRelative(ctx, spouseID, nodeID, Spouse)
			err = storage.setRelative(ctx, nodeID, spouseID, Spouse)
			if err != nil {
				customErr := fmt.Errorf("something went wrong while writing spouse, %v", err)

				utils.LogError(logger, customErr)
				log.Println(customErr)

				return err
			}
		}
	}

	return nil
}

func (storage *NodeStorage) setRelative(ctx context.Context, relativeID, nodeID uint, relationType string) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	_, err := storage.pool.Exec(ctx, SetRelativeQuery, relativeID, nodeID, relationType)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return err
	}

	return nil
}

func (storage *NodeStorage) writeAdditionalData(ctx context.Context, additional *models.AdditionDataList,
	nodeID uint) (*models.ReturningAdditionalData, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var (
		retData models.ReturningAdditionalData
		err     error
	)

	if additional.Birthdate != "" {
		retData.Birthdate, err = storage.updateBirthdate(ctx, additional.Birthdate, nodeID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while writing birthdate, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}
	} else {
		retData.Birthdate = nil
	}

	if additional.Deathdate != "" {
		retData.Deathdate, err = storage.updateDeathdate(ctx, additional.Deathdate, nodeID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while writing deathdate, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}
	} else {
		retData.Deathdate = nil
	}

	if additional.Description != "" {
		err = storage.updateDescription(ctx, additional.Description, nodeID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while writing description, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}
	}

	return &retData, nil
}

func (storage *NodeStorage) updateBirthdate(ctx context.Context, birthdate string,
	nodeID uint) (*pgtype.Date, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var writtenBirthdate *pgtype.Date

	birthdateToWrite, err := time.Parse("02/01/2006", birthdate)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while parsing date, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	line := storage.pool.QueryRow(ctx, UpdateBirthdateQuery, birthdateToWrite, nodeID)
	if err := line.Scan(&writtenBirthdate); err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return writtenBirthdate, nil
}

func (storage *NodeStorage) updateDeathdate(ctx context.Context, deathdate string,
	nodeID uint) (*pgtype.Date, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var writtenDeathdate *pgtype.Date

	deathdateToWrite, err := time.Parse("02/01/2006", deathdate)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while parsing date, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	line := storage.pool.QueryRow(ctx, UpdateDeathdateQuery, deathdateToWrite, nodeID)
	if err := line.Scan(&writtenDeathdate); err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return writtenDeathdate, nil
}

func (storage *NodeStorage) updateDescription(ctx context.Context, description string, nodeID uint) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var oldDescription string

	oldDescriptionLine := storage.pool.QueryRow(ctx, GetDescriptionQuery, nodeID)
	if oldDescriptionErr := oldDescriptionLine.Scan(&oldDescription); oldDescriptionErr != nil {
		_, err := storage.pool.Exec(ctx, UpdateDescriptionQuery, description, nodeID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while executing query, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return err
		}
	}

	_, err := storage.pool.Exec(ctx, InsertDescriptionQuery, description, nodeID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return err
	}

	return nil
}

func (storage *NodeStorage) updateGender(ctx context.Context, gender string, nodeID uint) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	_, err := storage.pool.Exec(ctx, UpdateGenderQuery, gender, nodeID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return err
	}

	return nil
}

func (storage *NodeStorage) updateName(ctx context.Context, name string, nodeID uint) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	_, err := storage.pool.Exec(ctx, UpdateNameQuery, name, nodeID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return err
	}

	return nil
}

func (storage *NodeStorage) createLayerIfNotExists(ctx context.Context, treeID uint, number int) (uint, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var layerID uint

	layerID, err := storage.getLayer(ctx, treeID, number)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		customErr := fmt.Errorf("something went wrong while getting layer ID, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return 0, err
	} else if err != nil {
		line := storage.pool.QueryRow(ctx, CreateLayerQuery, treeID, number)
		if err := line.Scan(&layerID); err != nil {
			customErr := fmt.Errorf("something went wrong while executing query, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return 0, err
		}
	}

	return layerID, nil
}
