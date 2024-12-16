package delivery

import (
	"encoding/json"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	proto "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (adminHandler *AdminHandler) Login(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	client := adminHandler.adminClient

	oldSession, _ := request.Cookie("admin_sid")
	if oldSession != nil {
		if oldUser, _ := client.GetAdminBySessionID(ctx,
			&proto.SessionRequest{SessionID: oldSession.Value}); oldUser != nil {
			log.Println("User already authorized", responses.StatusBadRequest)
			responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrAuthorized)

			return
		}
	}

	var requestData models.UserLoginRequest

	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	user, err := client.Login(ctx, &proto.LoginUserRequest{
		Email:    requestData.Email,
		Password: requestData.Password,
	})
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	newSession := adminHandler.createSession(user.SessionID)
	http.SetCookie(writer, newSession)

	responses.SendOkResponse(writer, &models.AdminResponse{
		ID:     uint(user.ID),
		Email:  user.Email,
		IsAuth: true,
	})
}

func (adminHandler *AdminHandler) Logout(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	session, _ := request.Cookie("admin_sid")

	client := adminHandler.adminClient

	_, err := client.Logout(ctx, &proto.SessionRequest{SessionID: session.Value})
	if err != nil {
		log.Println("User not authorized", responses.StatusUnauthorized)
		responses.SendErrResponse(writer, logger, responses.StatusUnauthorized, responses.ErrNotAuthorized)

		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(writer, session)

	responses.SendOkResponse(writer, &models.AdminResponse{IsAuth: false})
}

func (adminHandler *AdminHandler) EditUserPassword(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	var requestData models.EditPasswordRequest

	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	client := adminHandler.adminClient

	user, err := client.EditPassword(ctx, &proto.EditPasswordRequest{
		ID:       uint32(requestData.ID),
		Password: requestData.Password,
	})
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, &models.UserForAdminResponse{
		ID:           uint(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
}

func (adminHandler *AdminHandler) GetUsersList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	client := adminHandler.adminClient

	list, err := client.GetUsersList(ctx, &emptypb.Empty{})
	if err != nil {
		log.Println(responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	var users models.UsersList

	for _, user := range list.Users {
		users = append(users, models.UserForAdminResponse{
			ID:           uint(user.ID),
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
		})
	}

	responses.SendOkResponse(writer, &users)
}

func (adminHandler *AdminHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	var requestData models.UserLoginRequest

	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	client := adminHandler.adminClient

	user, err := client.CreateUser(ctx, &proto.LoginUserRequest{
		Email:    requestData.Email,
		Password: requestData.Password,
	})
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, &models.UserForAdminResponse{
		ID:           uint(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
}

func (adminHandler *AdminHandler) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	client := adminHandler.adminClient

	resp, err := client.DeleteUser(ctx, &proto.DeleteUserRequest{ID: uint32(id)})
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, &models.SuccessResponse{Success: resp.GetSuccess()})
}

func (adminHandler *AdminHandler) CheckAuth(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	session, _ := request.Cookie("admin_sid")
	if session == nil {
		responses.SendOkResponse(writer, &models.AdminResponse{IsAuth: false})

		return
	}

	client := adminHandler.adminClient
	admin, _ := client.GetAdminBySessionID(ctx, &proto.SessionRequest{SessionID: session.Value})
	if admin == nil {
		responses.SendOkResponse(writer, &models.AdminResponse{IsAuth: false})

		return
	}

	responses.SendOkResponse(writer, &models.AdminResponse{
		ID:     uint(admin.ID),
		Email:  admin.Email,
		IsAuth: true,
	})
}
