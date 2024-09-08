package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
)

func (storage *ProfileStorage) CreateProfile(ctx context.Context, userID uint) (*models.Profile, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var profile models.Profile

	line := storage.pool.QueryRow(ctx, CreateProfileQuery, userID)
	if err := line.Scan(&profile.ID, &profile.UserID); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)
		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	return &profile, nil
}

func (storage *ProfileStorage) UpdateProfile(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	return nil, nil
}
