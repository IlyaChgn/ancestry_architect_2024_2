package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	session "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/repository/session"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	profileusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/usecases"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandler struct {
	storage        usecases.AuthStorageInterface
	profileStorage profileusecases.ProfileStorageInterface
}

func NewAuthHandler(storage usecases.AuthStorageInterface,
	profileStorage profileusecases.ProfileStorageInterface) *AuthHandler {
	return &AuthHandler{
		storage:        storage,
		profileStorage: profileStorage,
	}
}

func createSession(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(session.SessionDuration),
		HttpOnly: true,
	}
}

func (authHandler *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	var requestData models.UserLoginRequest

	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	storage := authHandler.storage

	user, err := storage.GetUserByEmail(ctx, requestData.Email)
	if err != nil {
		log.Println("User with same email doesn`t exist", responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrWrongCredentials)

		return
	}

	if !utils.CheckPassword(requestData.Password, user.User.PasswordHash) {
		log.Println("Wrong password", responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrWrongCredentials)

		return
	}

	sessionID := uuid.NewString()
	err = storage.CreateSession(ctx, sessionID, user.User.ID)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	newSession := createSession(sessionID)
	http.SetCookie(writer, newSession)

	user.IsAuth = true
	responses.SendOkResponse(writer, user)
}

func (authHandler *AuthHandler) Logout(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	session, _ := request.Cookie("session_id")

	storage := authHandler.storage
	err := storage.RemoveSession(ctx, session.Value)

	if err != nil {
		log.Println("User not authorized", responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrNotAuthorized)

		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(writer, session)

	responses.SendOkResponse(writer, models.UserResponse{IsAuth: false})
}

func (authHandler *AuthHandler) Signup(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	var requestData models.UserSignupRequest

	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	storage := authHandler.storage

	user, errs := storage.CreateUser(ctx, requestData.Email, requestData.Password, requestData.PasswordRepeat)
	if errs != nil {
		log.Println(errs, responses.StatusBadRequest)
		responses.SendSeveralErrsResponse(writer, logger, responses.StatusBadRequest, errs)

		return
	}

	profileStorage := authHandler.profileStorage

	profile, err := profileStorage.CreateProfile(ctx, user.ID)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	sessionID := uuid.NewString()
	err = storage.CreateSession(ctx, sessionID, user.ID)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	newSession := createSession(sessionID)
	http.SetCookie(writer, newSession)

	responses.SendOkResponse(writer, models.UserResponse{
		User:    *user,
		IsAuth:  true,
		Name:    profile.Name,
		Surname: profile.Surname,
	})
}

func (authHandler *AuthHandler) CheckAuth(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	session, _ := request.Cookie("session_id")
	if session == nil {
		responses.SendOkResponse(writer, models.UserResponse{IsAuth: false})

		return
	}

	storage := authHandler.storage
	user, _ := storage.GetUserBySessionID(ctx, session.Value)

	if user == nil {
		responses.SendOkResponse(writer, models.UserResponse{IsAuth: false})

		return
	}

	user.IsAuth = true

	responses.SendOkResponse(writer, user)
}
