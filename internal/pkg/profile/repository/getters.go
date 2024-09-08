package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
)

func (storage *ProfileStorage) GetProfileByID(ctx context.Context, id uint) (*models.Profile, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var (
		profile     models.Profile
		nullProfile models.ProfileNullData
	)

	line := storage.pool.QueryRow(ctx, GetProfileByIDQuery, id)
	if err := line.Scan(&profile.ID, &profile.UserID, &nullProfile.Name, &nullProfile.Surname, &profile.Birthdate,
		&nullProfile.Gender, &nullProfile.AvatarPath); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)
		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	if nullProfile.Name != nil {
		profile.Name = *nullProfile.Name
	}

	if nullProfile.Surname != nil {
		profile.Surname = *nullProfile.Surname
	}

	if nullProfile.AvatarPath != nil {
		profile.AvatarPath = *nullProfile.AvatarPath
	}

	if nullProfile.Gender != nil {
		profile.Gender = *nullProfile.Gender
	}

	return &profile, nil
}
