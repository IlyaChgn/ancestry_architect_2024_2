package usecases

import (
	"context"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
)

type ProfileStorageInterface interface {
	GetProfileByID(ctx context.Context, id uint) (*models.Profile, error)

	CreateProfile(ctx context.Context, userID uint) (*models.Profile, error)
	UpdateProfile(ctx context.Context, profile *models.UpdateProfileRequest, userID uint) (*models.Profile, error)
}
