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

const (
	errUserAlreadyExists = "User with same emaul already exists"
)

func (storage *AuthStorage) createUser(ctx context.Context, tx pgx.Tx,
	email, passwordHash string) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	if _, err := tx.Exec(ctx, CreateUserQuery, email, passwordHash); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return err
	}

	return nil
}

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

	err := pgx.BeginFunc(ctx, storage.pool, func(tx pgx.Tx) error {
		err := storage.createUser(ctx, tx, email, utils.HashPassword(password))

		return err
	})

	if err != nil {
		utils.LogError(logger, fmt.Errorf("something went wrong while creating user, %v", err))

		return nil, []string{fmt.Sprintf("%v", err)}
	}

	user, _ := storage.GetUserByEmail(ctx, email)

	return user, nil
}

func (storage *AuthStorage) CreateSession(ctx context.Context, sessionID string, userID uint) error {
	return storage.manager.CreateSession(ctx, sessionID, userID)
}

func (storage *AuthStorage) RemoveSession(ctx context.Context, sessionID string) error {
	return storage.manager.RemoveSession(ctx, sessionID)
}
