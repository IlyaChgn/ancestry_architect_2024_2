package usecases

import (
	"context"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
)

type AuthStorageInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*models.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*models.UserResponse, error)
	GetUserBySessionID(ctx context.Context, sessionID string) (*models.UserResponse, error)

	CreateSession(ctx context.Context, sessionID string, userID uint) error
	RemoveSession(ctx context.Context, sessionID string) error
	CreateUser(ctx context.Context, email, password, passwordRepeat string) (*models.User, []string)
}
