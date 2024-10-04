package repository_test

import (
	"context"
	"testing"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/repository"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxIface(ctrl)

	email := "test@test.ru"
	name := ""
	surname := ""

	tests := []struct {
		name        string
		email       string
		expected    *models.UserResponse
		expectedErr bool
		setup       func()
	}{
		{
			name:  "successful case",
			email: email,
			expected: &models.UserResponse{
				User: models.User{
					ID:           1,
					Email:        email,
					PasswordHash: "hash",
				},
				Name:    "",
				Surname: "",
			},
			expectedErr: false,
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(uint(1), email, "hash", &name, &surname)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByEmailQuery, email).
					Return(pgxRows)
			},
		},
		{
			name:        "error case",
			email:       email,
			expected:    nil,
			expectedErr: true,
			setup: func() {
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByEmailQuery, email).
					Return(models.EmptyRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			storage := repository.NewAuthStorage(mockPool, nil)

			got, err := storage.GetUserByEmail(context.Background(), tt.email)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxIface(ctrl)

	id := uint(1)
	name := ""
	surname := ""

	tests := []struct {
		name        string
		id          uint
		expected    *models.UserResponse
		expectedErr bool
		setup       func()
	}{
		{
			name: "successful case",
			id:   id,
			expected: &models.UserResponse{
				User: models.User{
					ID:           1,
					Email:        "test@test.ru",
					PasswordHash: "hash",
				},
				Name:    name,
				Surname: surname,
			},
			expectedErr: false,
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(id, "test@test.ru", "hash", &name, &surname)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByIDQuery, id).
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
					QueryRow(context.Background(), repository.GetUserByIDQuery, id).
					Return(models.EmptyRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			storage := repository.NewAuthStorage(mockPool, nil)

			got, err := storage.GetUserByID(context.Background(), tt.id)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
