package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (storage *AuthStorage) getUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var user models.User

	line := tx.QueryRow(ctx, GerUserByEmailQuery, email)
	if err := line.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return &user, nil
}

func (storage *AuthStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var user *models.User

	err := pgx.BeginFunc(ctx, storage.pool, func(tx pgx.Tx) error {
		result, err := storage.getUserByEmail(ctx, tx, email)
		user = result

		return err
	})

	if err != nil {
		utils.LogError(logger, fmt.Errorf("something went wrong while getting user by email, %s", ErrUserNotExists))

		return nil, err
	}

	return user, nil
}

func (storage *AuthStorage) getUserByID(ctx context.Context, tx pgx.Tx, id uint) (*models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var user models.User

	line := tx.QueryRow(ctx, GerUserByIDQuery, id)
	if err := line.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)
		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return &user, nil
}

func (storage *AuthStorage) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var user *models.User

	err := pgx.BeginFunc(ctx, storage.pool, func(tx pgx.Tx) error {
		result, err := storage.getUserByID(ctx, tx, id)
		user = result

		return err
	})

	if err != nil {
		utils.LogError(logger, fmt.Errorf("something went wrong while getting user by ID, %s", ErrUserNotExists))

		return nil, err
	}

	return user, nil
}

func (storage *AuthStorage) GetUserBySessionID(ctx context.Context, sessionID string) (*models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	id, exists := storage.manager.GetUserID(ctx, sessionID)

	if !exists {
		utils.LogError(logger, fmt.Errorf("something went wrong while getting user by ID, %s", ErrUserNotExists))

		return nil, fmt.Errorf("something went wrong while getting user by ID, %s", ErrUserNotExists)
	}

	return storage.GetUserByID(ctx, id)
}
