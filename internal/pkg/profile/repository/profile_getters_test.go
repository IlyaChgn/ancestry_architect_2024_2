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
)

func TestGetProfileByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxIface(ctrl)

	id := uint(1)
	name := "name"
	surname := "surname"
	avatarPath := "avatarPath"
	gender := "gender"

	var nullProfile models.ProfileNullData

	tests := []struct {
		name        string
		id          uint
		expected    *models.Profile
		expectedErr bool
		setup       func()
	}{
		{
			name: "successful case",
			id:   id,
			expected: &models.Profile{
				ID:         id,
				UserID:     id,
				Name:       name,
				Surname:    surname,
				Birthdate:  pgtype.Date{},
				Gender:     gender,
				AvatarPath: avatarPath,
			},
			expectedErr: false,
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(id, id, &name, &surname, pgtype.Date{}, &gender, &avatarPath)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetProfileByIDQuery, id).
					Return(pgxRows)
			},
		},
		{
			name: "successful with null data",
			id:   id,
			expected: &models.Profile{
				ID:        id,
				UserID:    id,
				Birthdate: pgtype.Date{},
			},
			expectedErr: false,
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(id, id, nullProfile.Name, nullProfile.Surname, pgtype.Date{}, nullProfile.Gender, nullProfile.AvatarPath)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetProfileByIDQuery, id).
					Return(pgxRows)
			},
		},
		{
			name:        "error case",
			id:          id,
			expected:    nil,
			expectedErr: true,
			setup: func() {
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetProfileByIDQuery, id).
					Return(models.EmptyRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			storage := repository.NewProfileStorage(mockPool, "")

			got, err := storage.GetProfileByID(context.Background(), tt.id)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
