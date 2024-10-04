package repository_test

import (
	"context"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/repository"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateProfile(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxIface(ctrl)

	id := uint(1)

	tests := []struct {
		name        string
		userID      uint
		expected    *models.Profile
		expectedErr bool
		setup       func()
	}{
		{
			name:   "successful case",
			userID: id,
			expected: &models.Profile{
				ID:     id,
				UserID: id,
			},
			expectedErr: false,
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(id, id)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.CreateProfileQuery, id).
					Return(pgxRows)
			},
		},
		{
			name:        "error case",
			userID:      id,
			expectedErr: true,
			setup: func() {
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.CreateProfileQuery, id).
					Return(models.EmptyRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			storage := repository.NewProfileStorage(mockPool, "")

			got, err := storage.CreateProfile(context.Background(), tt.userID)
			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxIface(ctrl)

	name := "name"
	surname := "surname"
	gender := "gender"
	avatarPath := "avatarPath"
	birthdate := time.Now()

	tests := []struct {
		name        string
		userID      uint
		profile     *models.UpdateProfileRequest
		expected    *models.Profile
		expectedErr bool
		setup       func()
	}{
		{
			name:   "successful case",
			userID: uint(1),
			profile: &models.UpdateProfileRequest{
				Name:      name,
				Surname:   surname,
				Gender:    gender,
				Birthdate: birthdate,
			},
			expected: &models.Profile{
				ID:         uint(1),
				UserID:     uint(1),
				Name:       name,
				Surname:    surname,
				Birthdate:  pgtype.Date{},
				Gender:     gender,
				AvatarPath: avatarPath,
			},
			expectedErr: false,
			setup: func() {
				mockPool.EXPECT().Exec(context.Background(), repository.UpdateBirthdateQuery, birthdate, uint(1))
				mockPool.EXPECT().Exec(context.Background(), repository.UpdateGenderQuery, gender, uint(1))
				pgxRows := pgxpoolmock.NewRow(uint(1), uint(1), &name, &surname, pgtype.Date{}, &gender, &avatarPath)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.UpdateProfileQuery, name, surname, uint(1)).
					Return(pgxRows)
			},
		},
		{
			name:   "successful case without some fields",
			userID: uint(1),
			profile: &models.UpdateProfileRequest{
				Name:    name,
				Surname: surname,
			},
			expected: &models.Profile{
				ID:         uint(1),
				UserID:     uint(1),
				Name:       name,
				Surname:    surname,
				Birthdate:  pgtype.Date{},
				Gender:     gender,
				AvatarPath: avatarPath,
			},
			expectedErr: false,
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(uint(1), uint(1), &name, &surname, pgtype.Date{}, &gender, &avatarPath)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.UpdateProfileQuery, name, surname, uint(1)).
					Return(pgxRows)
			},
		},
		{
			name:   "error case",
			userID: uint(1),
			profile: &models.UpdateProfileRequest{
				Name:    name,
				Surname: surname,
			},
			expectedErr: true,
			setup: func() {
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.UpdateProfileQuery, name, surname, uint(1)).
					Return(models.EmptyRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			storage := repository.NewProfileStorage(mockPool, "")

			got, err := storage.UpdateProfile(context.Background(), tt.profile, tt.userID)
			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
