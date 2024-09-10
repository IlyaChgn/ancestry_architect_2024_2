package repository

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

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

func (storage *ProfileStorage) UpdateProfile(ctx context.Context, profile *models.UpdateProfileRequest,
	userID uint) (*models.Profile, error) {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	if profile.Avatar != nil {
		err := storage.updateProfileAvatar(ctx, profile.Avatar, userID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while updating avatar, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}
	}

	if !profile.Birthdate.IsZero() {
		err := storage.updateProfileBirthdate(ctx, profile.Birthdate, userID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while updating birthdate, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}
	}

	if profile.Gender != "" {
		err := storage.updateProfileGender(ctx, profile.Gender, userID)
		if err != nil {
			customErr := fmt.Errorf("something went wrong while updating gender, %v", err)

			utils.LogError(logger, customErr)
			log.Println(customErr)

			return nil, err
		}
	}

	var (
		newProfile  models.Profile
		nullProfile models.ProfileNullData
	)

	line := storage.pool.QueryRow(ctx, UpdateProfileQuery, profile.Name, profile.Surname, userID)
	if err := line.Scan(&newProfile.ID, &newProfile.UserID, &nullProfile.Name, &nullProfile.Surname,
		&newProfile.Birthdate, &nullProfile.Gender, &nullProfile.AvatarPath); err != nil {
		customErr := fmt.Errorf("something went wrong while scanning line, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return nil, err
	}

	if nullProfile.Name != nil {
		newProfile.Name = *nullProfile.Name
	}

	if nullProfile.Surname != nil {
		newProfile.Surname = *nullProfile.Surname
	}

	if nullProfile.AvatarPath != nil {
		newProfile.AvatarPath = *nullProfile.AvatarPath
	}

	if nullProfile.Gender != nil {
		newProfile.Gender = *nullProfile.Gender
	}

	return &newProfile, nil
}

func (storage *ProfileStorage) updateProfileAvatar(ctx context.Context, avatar *multipart.FileHeader,
	userID uint) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	var oldPath string

	oldPathLine := storage.pool.QueryRow(ctx, GetAvatarQuery, userID)
	if oldPathErr := oldPathLine.Scan(&oldPath); oldPathErr == nil {
		os.Remove(oldPath)
	}

	newPath, err := utils.WriteFile(avatar, "avatar")
	if err != nil {
		customErr := fmt.Errorf("something went wrong while writing avatar, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return err
	}

	_, err = storage.pool.Exec(ctx, UpdateAvatarQuery, newPath, userID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return err
	}

	return nil
}

func (storage *ProfileStorage) updateProfileBirthdate(ctx context.Context, birthdate time.Time, userID uint) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	_, err := storage.pool.Exec(ctx, UpdateBirthdateQuery, birthdate, userID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return err
	}

	return nil
}

func (storage *ProfileStorage) updateProfileGender(ctx context.Context, gender string, userID uint) error {
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("storage", utils.GetFunctionName()))

	_, err := storage.pool.Exec(ctx, UpdateGenderQuery, gender, userID)
	if err != nil {
		customErr := fmt.Errorf("something went wrong while executing query, %v", err)

		utils.LogError(logger, customErr)
		log.Println(customErr)

		return err
	}

	return nil
}
