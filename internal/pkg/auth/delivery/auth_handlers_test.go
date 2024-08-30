package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/delivery"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	mock_usecases "github.com/IlyaChgn/ancestry_architect_2024_2/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthStorage := mock_usecases.NewMockAuthStorageInterface(ctrl)
	hash := utils.HashPassword("password")

	tests := []struct {
		name           string
		requestBody    []byte
		expectedCode   int
		expectedStatus string
		setup          func()
	}{
		{
			name:         "successful login",
			requestBody:  []byte(`{"email": "", "password": "password"}`),
			expectedCode: responses.StatusOk,
			setup: func() {
				mockAuthStorage.EXPECT().GetUserByEmail(context.Background(), gomock.Any()).Return(&models.User{
					ID:           1,
					PasswordHash: hash,
				}, nil)
				mockAuthStorage.EXPECT().CreateSession(context.Background(), gomock.Any(), uint(1)).Return(nil)
			},
		},
		{
			name:           "test wrong password",
			requestBody:    []byte(`{"email": "", "password": "password1"}`),
			expectedCode:   responses.StatusBadRequest,
			expectedStatus: responses.ErrWrongCredentials,
			setup: func() {
				mockAuthStorage.EXPECT().GetUserByEmail(context.Background(), gomock.Any()).Return(&models.User{
					ID:           1,
					PasswordHash: hash,
				}, nil)
			},
		},
		{
			name:           "test internal server error",
			requestBody:    []byte(`{"email": "", "password": "password"}`),
			expectedCode:   responses.StatusInternalServerError,
			expectedStatus: responses.ErrInternalServer,
			setup: func() {
				mockAuthStorage.EXPECT().GetUserByEmail(context.Background(), gomock.Any()).Return(&models.User{
					ID:           1,
					PasswordHash: hash,
				}, nil)
				mockAuthStorage.EXPECT().CreateSession(context.Background(), gomock.Any(), uint(1)).
					Return(errors.New("err"))
			},
		},
		{
			name:           "test user does not exist",
			requestBody:    []byte(`{"email": "", "password": "password"}`),
			expectedCode:   responses.StatusBadRequest,
			expectedStatus: responses.ErrWrongCredentials,
			setup: func() {
				mockAuthStorage.EXPECT().GetUserByEmail(context.Background(), gomock.Any()).Return(nil,
					errors.New("User with same email doesn`t exist"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			authHandler := delivery.NewAuthHandler(mockAuthStorage)

			request := httptest.NewRequest("POST", "/login", bytes.NewBuffer(tt.requestBody))

			writer := httptest.NewRecorder()
			handle := http.HandlerFunc(authHandler.Login)

			handle.ServeHTTP(writer, request)

			if tt.expectedStatus != "" {
				var errStructure *models.ErrResponse

				err := json.NewDecoder(writer.Body).Decode(&errStructure)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.expectedCode, writer.Code)
				assert.Equal(t, errStructure.Status, tt.expectedStatus)
			} else {
				assert.Equal(t, responses.StatusOk, writer.Code)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthStorage := mock_usecases.NewMockAuthStorageInterface(ctrl)
	cookieName := "session_id"

	tests := []struct {
		name           string
		expectedCode   int
		expectedStatus string
		setup          func()
	}{
		{
			name:         "successful logout",
			expectedCode: responses.StatusOk,
			setup: func() {
				mockAuthStorage.EXPECT().RemoveSession(context.Background(), gomock.Any()).Return(nil)
			},
		},
		{
			name:           "bad logout",
			expectedCode:   responses.StatusBadRequest,
			expectedStatus: responses.ErrNotAuthorized,
			setup: func() {
				mockAuthStorage.EXPECT().RemoveSession(context.Background(), gomock.Any()).
					Return(errors.New("User not authorized"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			authHandler := delivery.NewAuthHandler(mockAuthStorage)

			request := httptest.NewRequest("POST", "/logout", nil)

			request.AddCookie(&http.Cookie{Name: cookieName})

			writer := httptest.NewRecorder()
			handle := http.HandlerFunc(authHandler.Logout)

			handle.ServeHTTP(writer, request)

			if tt.expectedStatus != "" {
				var errStructure *models.ErrResponse

				err := json.NewDecoder(writer.Body).Decode(&errStructure)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.expectedCode, writer.Code)
				assert.Equal(t, errStructure.Status, tt.expectedStatus)
			} else {
				assert.Equal(t, responses.StatusOk, writer.Code)
			}
		})
	}
}

func TestSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthStorage := mock_usecases.NewMockAuthStorageInterface(ctrl)

	tests := []struct {
		name           string
		requestBody    []byte
		expectedCode   int
		expectedStatus string
		hasSeveralErrs bool
		severalErrs    []string
		setup          func()
	}{
		{
			name:         "successful signup",
			requestBody:  []byte(`{"email": "", "password": "password"}`),
			expectedCode: responses.StatusOk,
			setup: func() {
				mockAuthStorage.EXPECT().CreateUser(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.User{
						ID: 1,
					}, nil)
				mockAuthStorage.EXPECT().CreateSession(context.Background(), gomock.Any(), uint(1)).Return(nil)
			},
		},
		{
			name:           "test wrong credentials",
			requestBody:    []byte(`{"email": "", "password": "password"}`),
			expectedCode:   responses.StatusBadRequest,
			hasSeveralErrs: true,
			severalErrs:    []string{"Passwords do not match"},
			setup: func() {
				mockAuthStorage.EXPECT().CreateUser(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, []string{"Passwords do not match"})
			},
		},
		{
			name:           "test internal server error",
			requestBody:    []byte(`{"email": "", "password": "password"}`),
			expectedCode:   responses.StatusInternalServerError,
			expectedStatus: responses.ErrInternalServer,
			setup: func() {
				mockAuthStorage.EXPECT().CreateUser(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.User{
						ID: 1,
					}, nil)

				mockAuthStorage.EXPECT().CreateSession(context.Background(), gomock.Any(), uint(1)).
					Return(errors.New("err"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			authHandler := delivery.NewAuthHandler(mockAuthStorage)

			request := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(tt.requestBody))

			writer := httptest.NewRecorder()
			handle := http.HandlerFunc(authHandler.Signup)

			handle.ServeHTTP(writer, request)

			if tt.hasSeveralErrs {
				var errStructure *models.SeveralErrsResponse

				err := json.NewDecoder(writer.Body).Decode(&errStructure)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.expectedCode, writer.Code)
				assert.Equal(t, errStructure.Errors, tt.severalErrs)
			} else if tt.expectedStatus != "" {
				var errStructure *models.ErrResponse

				err := json.NewDecoder(writer.Body).Decode(&errStructure)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.expectedCode, writer.Code)
				assert.Equal(t, errStructure.Status, tt.expectedStatus)
			} else {
				assert.Equal(t, responses.StatusOk, writer.Code)
			}
		})
	}
}

func TestCheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthStorage := mock_usecases.NewMockAuthStorageInterface(ctrl)
	cookieName := "session_id"

	tests := []struct {
		name       string
		needCookie bool
		setup      func()
	}{
		{
			name:       "test is auth",
			needCookie: true,
			setup: func() {
				mockAuthStorage.EXPECT().GetUserBySessionID(context.Background(), gomock.Any()).Return(&models.User{
					ID: 1,
				}, nil)
			},
		},
		{
			name:       "test not auth with cookie",
			needCookie: true,
			setup: func() {
				mockAuthStorage.EXPECT().GetUserBySessionID(context.Background(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:       "test not auth without cookie",
			needCookie: false,
			setup:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			authHandler := delivery.NewAuthHandler(mockAuthStorage)

			request := httptest.NewRequest("GET", "/check_auth", nil)

			if tt.needCookie {
				request.AddCookie(&http.Cookie{Name: cookieName})
			}

			writer := httptest.NewRecorder()
			handle := http.HandlerFunc(authHandler.CheckAuth)

			handle.ServeHTTP(writer, request)

			assert.Equal(t, responses.StatusOk, writer.Code)
		})
	}
}
