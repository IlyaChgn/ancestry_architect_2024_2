package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
)

const (
	errUserAlreadyExists = "User with same email already exists"
)

func (storage *AuthStorage) CreateUser(ctx context.Context, email, password,
	passwordRepeat string) (*models.User, []string) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	oldUser, _ := storage.GetUserByEmail(ctx, email)
	if oldUser != nil {
		return nil, []string{errUserAlreadyExists}
	}

	errs := utils.Validate(email, password, passwordRepeat)
	if errs != nil {
		return nil, errs
	}

	var user models.User

	line := storage.pool.QueryRow(ctx, CreateUserQuery, email, utils.HashPassword(password))

	if err := line.Scan(&user.ID, &user.Email); err != nil {
		customErr := fmt.Errorf("something went wrong while creating user, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, []string{fmt.Sprintf("%v", customErr)}
	}

	return &user, nil
}

func (storage *AuthStorage) CreateSession(ctx context.Context, sessionID string, userID uint) error {
	return storage.manager.CreateSession(ctx, sessionID, userID)
}

func (storage *AuthStorage) RemoveSession(ctx context.Context, sessionID string) error {
	return storage.manager.RemoveSession(ctx, sessionID)
}
