package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
)

func (storage *AuthStorage) GetUserByEmail(ctx context.Context, email string) (*models.UserResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var (
		user        models.User
		nullProfile models.ProfileNullData
	)

	line := storage.pool.QueryRow(ctx, GetUserByEmailQuery, email)
	if err := line.Scan(&user.ID, &user.Email, &user.PasswordHash, &nullProfile.Name,
		&nullProfile.Surname); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	var name string
	if nullProfile.Name != nil {
		name = *nullProfile.Name
	}

	var surname string
	if nullProfile.Surname != nil {
		surname = *nullProfile.Surname
	}

	return &models.UserResponse{
		User:    user,
		Name:    name,
		Surname: surname,
	}, nil
}

func (storage *AuthStorage) GetUserByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var (
		user        models.User
		nullProfile models.ProfileNullData
	)

	line := storage.pool.QueryRow(ctx, GetUserByIDQuery, id)
	if err := line.Scan(&user.ID, &user.Email, &user.PasswordHash, &nullProfile.Name,
		&nullProfile.Surname, &nullProfile.AvatarPath); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	var name string
	if nullProfile.Name != nil {
		name = *nullProfile.Name
	}

	var surname string
	if nullProfile.Surname != nil {
		surname = *nullProfile.Surname
	}

	var avatarPath string
	if nullProfile.AvatarPath != nil {
		avatarPath = *nullProfile.AvatarPath
	}

	return &models.UserResponse{
		User:       user,
		Name:       name,
		Surname:    surname,
		AvatarPath: avatarPath,
	}, nil
}

func (storage *AuthStorage) GetUserBySessionID(ctx context.Context, sessionID string) (*models.UserResponse, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	id, exists := storage.manager.GetUserID(ctx, sessionID)

	if !exists {
		utils.LogError(logger, fmt.Errorf("something went wrong while getting user by session ID, %s",
			ErrUserNotExists))

		return nil, fmt.Errorf("something went wrong while getting user by session ID, %s", ErrUserNotExists)
	}

	return storage.GetUserByID(ctx, id)
}
