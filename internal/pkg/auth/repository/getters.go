package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
)

func (storage *AuthStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var user models.User

	line := storage.pool.QueryRow(ctx, GetUserByEmailQuery, email)
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

	var user models.User

	line := storage.pool.QueryRow(ctx, GetUserByIDQuery, id)
	if err := line.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)
		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return &user, nil
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
