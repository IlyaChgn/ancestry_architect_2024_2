package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/delivery"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	mock_usecases "github.com/IlyaChgn/ancestry_architect_2024_2/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProfile(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthStorage := mock_usecases.NewMockAuthStorageInterface(ctrl)
	mockProfileStorage := mock_usecases.NewMockProfileStorageInterface(ctrl)

	tests := []struct {
		name           string
		expectedCode   int
		expectedStatus string
		setup          func()
	}{
		{
			name:         "successful get profile",
			expectedCode: responses.StatusOk,
			setup: func() {
				mockProfileStorage.EXPECT().GetProfileByID(gomock.Any(), gomock.Any()).
					Return(&models.Profile{
						Birthdate: pgtype.Date{Status: pgtype.Null},
					}, nil)
			},
		},
		{
			name:           "test bad request",
			expectedCode:   responses.StatusBadRequest,
			expectedStatus: responses.ErrBadRequest,
			setup: func() {
				mockProfileStorage.EXPECT().GetProfileByID(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("test"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			profileHandler := delivery.NewProfileHandler(mockProfileStorage, mockAuthStorage)

			request := httptest.NewRequest("GET", "/profile/1", nil)
			request = mux.SetURLVars(request, map[string]string{"id": "1"})

			writer := httptest.NewRecorder()
			handle := http.HandlerFunc(profileHandler.GetProfile)

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

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthStorage := mock_usecases.NewMockAuthStorageInterface(ctrl)
	mockProfileStorage := mock_usecases.NewMockProfileStorageInterface(ctrl)
	cookieName := "session_id"

	tests := []struct {
		name           string
		formData       map[string]string
		expectedCode   int
		expectedStatus string
		hasSeveralErrs bool
		severalErrs    []string
		setup          func()
	}{
		{
			name: "successful update profile",
			formData: map[string]string{
				"email":     "test@test.com",
				"birthdate": "30/09/2004",
			},
			expectedCode: responses.StatusOk,
			setup: func() {
				mockAuthStorage.EXPECT().GetUserBySessionID(context.Background(), gomock.Any()).
					Return(&models.UserResponse{}, nil)
				mockAuthStorage.EXPECT().UpdateEmail(context.Background(), gomock.Any(), gomock.Any()).
					Return(&models.User{}, nil)
				mockProfileStorage.EXPECT().UpdateProfile(context.Background(), gomock.Any(), gomock.Any()).
					Return(&models.Profile{
						Birthdate: pgtype.Date{Status: pgtype.Null},
					}, nil)
			},
		},
		{
			name:         "successful update profile without birthdate and email",
			expectedCode: responses.StatusOk,
			setup: func() {
				mockAuthStorage.EXPECT().GetUserBySessionID(context.Background(), gomock.Any()).
					Return(&models.UserResponse{}, nil)
				mockProfileStorage.EXPECT().UpdateProfile(context.Background(), gomock.Any(), gomock.Any()).
					Return(&models.Profile{
						Birthdate: pgtype.Date{Status: pgtype.Null},
					}, nil)
			},
		},
		{
			name:           "test internal server error",
			expectedCode:   responses.StatusInternalServerError,
			expectedStatus: responses.ErrInternalServer,
			setup: func() {
				mockAuthStorage.EXPECT().GetUserBySessionID(context.Background(), gomock.Any()).
					Return(nil, errors.New("test error"))
			},
		},
		{
			name:           "test wrong date format",
			expectedCode:   responses.StatusBadRequest,
			expectedStatus: responses.ErrBadRequest,
			formData: map[string]string{
				"birthdate": "30.09.2004",
			},
			setup: func() {
				mockAuthStorage.EXPECT().GetUserBySessionID(context.Background(), gomock.Any()).
					Return(&models.UserResponse{}, nil)
			},
		},
		{
			name:           "test wrong email",
			expectedCode:   responses.StatusBadRequest,
			expectedStatus: responses.ErrBadRequest,
			formData: map[string]string{
				"email": "testtest.com",
			},
			hasSeveralErrs: true,
			severalErrs:    []string{"Wrong email format"},
			setup: func() {
				mockAuthStorage.EXPECT().GetUserBySessionID(context.Background(), gomock.Any()).
					Return(&models.UserResponse{}, nil)
				mockAuthStorage.EXPECT().UpdateEmail(context.Background(), gomock.Any(), gomock.Any()).
					Return(nil, []string{"Wrong email format"})
			},
		},
		{
			name:           "test bad request",
			expectedCode:   responses.StatusBadRequest,
			expectedStatus: responses.ErrBadRequest,
			setup: func() {
				mockAuthStorage.EXPECT().GetUserBySessionID(context.Background(), gomock.Any()).
					Return(&models.UserResponse{}, nil)
				mockProfileStorage.EXPECT().UpdateProfile(context.Background(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("test error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			profileHandler := delivery.NewProfileHandler(mockProfileStorage, mockAuthStorage)

			body := new(bytes.Buffer)
			formWriter := multipart.NewWriter(body)
			for key, value := range tt.formData {
				part, err := formWriter.CreateFormField(key)
				if err != nil {
					t.Fatalf("failed to create form field: %v", err)
				}
				_, err = part.Write([]byte(value))
				if err != nil {
					t.Fatalf("failed to write form field: %v", err)
				}
			}

			err := formWriter.Close()
			if err != nil {
				t.Fatalf("failed to close multipart writer: %v", err)
			}

			request, err := http.NewRequest("POST", "/profile/edit", body)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			request.Header.Set("Content-Type", formWriter.FormDataContentType())

			cookie := &http.Cookie{Name: cookieName, Value: "some_session_id"}
			request.AddCookie(cookie)

			writer := httptest.NewRecorder()
			handle := http.HandlerFunc(profileHandler.UpdateProfile)

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
