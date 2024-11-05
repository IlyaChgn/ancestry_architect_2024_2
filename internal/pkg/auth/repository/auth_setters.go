package repository

import (
	"context"
	"fmt"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	repository "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/session"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
	"log"
)

func (storage *AuthStorage) CreateUser(ctx context.Context, email, password,
	passwordRepeat string) (*models.User, []string) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	oldUser, _ := storage.GetUserByEmail(ctx, email)
	if oldUser != nil {
		return nil, []string{ErrUserAlreadyExists}
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
	return storage.manager.CreateSession(ctx, sessionID, userID, repository.UserSessionDuration)
}

func (storage *AuthStorage) RemoveSession(ctx context.Context, sessionID string) error {
	return storage.manager.RemoveSession(ctx, sessionID)
}

func (storage *AuthStorage) UpdateEmail(ctx context.Context, email string, userID uint) (*models.User, []string) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var errs []string

	oldUser, _ := storage.GetUserByEmail(ctx, email)
	if oldUser != nil && oldUser.User.ID != userID {
		errs = append(errs, ErrUserAlreadyExists)
	}

	if !utils.ValidateEmail(email) {
		errs = append(errs, ErrWrongEmailFormat)
	}

	if errs != nil {
		return nil, errs
	}

	var user models.User

	line := storage.pool.QueryRow(ctx, UpdateEmailQuery, email, userID)
	if err := line.Scan(&user.ID, &user.Email); err != nil {
		customErr := fmt.Errorf("something went wrong while creating user, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, []string{fmt.Sprintf("%v", customErr)}
	}

	return &user, nil
}
