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

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxIface(ctrl)

	email := "test@test.ru"
	name := ""
	surname := ""

	tests := []struct {
		name            string
		user            *models.UserSignupRequest
		expected        *models.User
		expectedErr     bool
		expectedErrList []string
		setup           func()
	}{
		{
			name: "successful case",
			user: &models.UserSignupRequest{
				UserLoginRequest: models.UserLoginRequest{
					Email:    email,
					Password: "password?",
				},
				PasswordRepeat: "password?",
			},
			expected: &models.User{
				ID:    1,
				Email: email,
			},
			expectedErr: false,
			setup: func() {
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByEmailQuery, email).
					Return(models.EmptyRow{})
				pgxRows := pgxpoolmock.NewRow(uint(1), email)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.CreateUserQuery, email, gomock.Any()).
					Return(pgxRows)
			},
		},
		{
			name: "error case 1",
			user: &models.UserSignupRequest{
				UserLoginRequest: models.UserLoginRequest{
					Email:    email,
					Password: "password?",
				},
				PasswordRepeat: "password?",
			},
			expected:        nil,
			expectedErr:     true,
			expectedErrList: []string{"User with same email already exists"},
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(uint(1), email, "", &name, &surname)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByEmailQuery, email).
					Return(pgxRows)
			},
		},
		{
			name: "error case 2",
			user: &models.UserSignupRequest{
				UserLoginRequest: models.UserLoginRequest{
					Email:    email,
					Password: "password?",
				},
				PasswordRepeat: "password?1",
			},
			expected:        nil,
			expectedErr:     true,
			expectedErrList: []string{"Passwords do not match"},
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

			got, err := storage.CreateUser(context.Background(), tt.user.Email, tt.user.Password,
				tt.user.PasswordRepeat)
			if tt.expectedErr {
				assert.Equal(t, tt.expectedErrList, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestUpdateEmail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxIface(ctrl)

	email := "test@test.ru"
	name := ""
	surname := ""

	tests := []struct {
		name            string
		user            *models.User
		expected        *models.User
		expectedErr     bool
		expectedErrList []string
		setup           func()
	}{
		{
			name: "successful test",
			user: &models.User{
				ID:    1,
				Email: email,
			},
			expected: &models.User{
				ID:    1,
				Email: email,
			},
			expectedErr: false,
			setup: func() {
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByEmailQuery, email).
					Return(models.EmptyRow{})
				pgxRows := pgxpoolmock.NewRow(uint(1), email)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.UpdateEmailQuery, email, uint(1)).
					Return(pgxRows)
			},
		},
		{
			name: "test user already exists",
			user: &models.User{
				ID:    1,
				Email: email,
			},
			expected: &models.User{
				ID:    1,
				Email: email,
			},
			expectedErr:     true,
			expectedErrList: []string{"User with same email already exists"},
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(uint(2), email, "", &name, &surname)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByEmailQuery, email).
					Return(pgxRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			storage := repository.NewAuthStorage(mockPool, nil)

			got, err := storage.UpdateEmail(context.Background(), tt.user.Email, tt.user.ID)
			if tt.expectedErr {
				assert.Equal(t, tt.expectedErrList, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
