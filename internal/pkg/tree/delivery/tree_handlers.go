package delivery

import (
	"encoding/json"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	authusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/tree/usecases"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

type TreeHandler struct {
	storage     usecases.TreeStorageInterface
	authStorage authusecases.AuthStorageInterface
}

func NewTreeHandler(storage usecases.TreeStorageInterface, authStorage authusecases.AuthStorageInterface) *TreeHandler {
	return &TreeHandler{
		storage:     storage,
		authStorage: authStorage,
	}
}

func (treeHandler *TreeHandler) GetTree(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	vars := mux.Vars(request)

	treeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	session, _ := request.Cookie("session_id")
	user, _ := treeHandler.authStorage.GetUserBySessionID(ctx, session.Value)

	storage := treeHandler.storage

	hasPermission, err := storage.CheckPermission(ctx, uint(treeID), user.User.ID)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	if !hasPermission {
		log.Println(responses.ErrForbidden, responses.StatusForbidden)
		responses.SendErrResponse(writer, logger, responses.StatusForbidden, responses.ErrForbidden)

		return
	}

	responses.SendOkResponse(writer, "as")
}

func (treeHandler *TreeHandler) GetCreatedTreesList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	session, _ := request.Cookie("session_id")
	user, _ := treeHandler.authStorage.GetUserBySessionID(ctx, session.Value)

	storage := treeHandler.storage

	treesList, err := storage.GetCreatedTrees(ctx, user.User.ID)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, treesList)
}

func (treeHandler *TreeHandler) GetAvailableTreesList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	session, _ := request.Cookie("session_id")
	user, _ := treeHandler.authStorage.GetUserBySessionID(ctx, session.Value)

	storage := treeHandler.storage

	treesList, err := storage.GetAvailableTrees(ctx, user.User.ID)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, treesList)
}

func (treeHandler *TreeHandler) CreateTree(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	var requestData models.CreateTreeRequest
	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	session, _ := request.Cookie("session_id")
	user, _ := treeHandler.authStorage.GetUserBySessionID(ctx, session.Value)

	storage := treeHandler.storage

	tree, err := storage.CreateTree(ctx, user.User.ID, requestData.Name)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, tree)
}

func (treeHandler *TreeHandler) AddPermission(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	var requestData models.AddPermissionRequest

	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	storage := treeHandler.storage

	err = storage.AddPermission(ctx, requestData.UserID, requestData.TreeID)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, &models.SuccessResponse{Success: true})
}
