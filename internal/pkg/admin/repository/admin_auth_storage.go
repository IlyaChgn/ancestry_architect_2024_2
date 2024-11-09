package repository

import (
	"context"
	"fmt"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/session"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
	"log"
)

func (storage *AdminStorage) GetAdminByEmail(ctx context.Context, email string) (*models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var user models.User

	line := storage.pool.QueryRow(ctx, GetAdminByEmailQuery, email)
	if err := line.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return &user, nil
}

func (storage *AdminStorage) GetAdminBySessionID(ctx context.Context, sessionID string) (*models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	id, exists := storage.manager.GetUserID(ctx, sessionID)
	if !exists {
		customErr := fmt.Errorf("something went wrong while getting admin by sessiion ID, user not found")

		log.Println(customErr)
		utils.LogError(logger, customErr)

		return nil, customErr
	}

	return storage.getUserByID(ctx, id)
}

func (storage *AdminStorage) UpdatePassword(ctx context.Context, id uint, password string) (*models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var user models.User

	line := storage.pool.QueryRow(ctx, EditUserPasswordQuery, id, utils.HashPassword(password))
	if err := line.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		customErr := fmt.Errorf("something went wrong while editing user password, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, customErr
	}

	return &user, nil
}

func (storage *AdminStorage) GetUsersList(ctx context.Context) (*[]models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var users []models.User

	rows, err := storage.pool.Query(ctx, GetUsersListQuery)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while getting users list, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, customErr
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
			customErr := fmt.Errorf("something went wrong while scanning row, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, customErr
		}

		users = append(users, user)
	}

	return &users, nil
}

func (storage *AdminStorage) getUserByID(ctx context.Context, id uint) (*models.User, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var user models.User

	line := storage.pool.QueryRow(ctx, GetAdminByIDQuery, id)
	if err := line.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return &user, nil
}

func (storage *AdminStorage) CreateSession(ctx context.Context, sessionID string, userID uint) error {
	return storage.manager.CreateSession(ctx, sessionID, userID, session.AdminSessionDuration)
}

func (storage *AdminStorage) RemoveSession(ctx context.Context, sessionID string) error {
	return storage.manager.RemoveSession(ctx, sessionID)
}
