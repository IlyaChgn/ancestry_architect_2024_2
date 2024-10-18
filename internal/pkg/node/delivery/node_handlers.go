package delivery

import (
	"encoding/json"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	authusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/node/usecases"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	treeusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/tree/usecases"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

type NodeHandler struct {
	storage     usecases.NodeStorageInterface
	treeStorage treeusecases.TreeStorageInterface
	authStorage authusecases.AuthStorageInterface
}

func NewNodeHandler(storage usecases.NodeStorageInterface, treeStorage treeusecases.TreeStorageInterface,
	authStorage authusecases.AuthStorageInterface) *NodeHandler {
	return &NodeHandler{
		storage:     storage,
		treeStorage: treeStorage,
		authStorage: authStorage,
	}
}

func (nodeHandler *NodeHandler) GetDescription(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	vars := mux.Vars(request)
	nodeID, _ := strconv.Atoi(vars["id"])

	storage := nodeHandler.storage

	description, err := storage.GetDescription(ctx, uint(nodeID))
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, description)
}

func (nodeHandler *NodeHandler) CreateNode(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	authStorage := nodeHandler.authStorage
	session, _ := request.Cookie("session_id")

	user, err := authStorage.GetUserBySessionID(ctx, session.Value)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	var requestData models.CreateNodeRequest

	err = json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	hasPermission, err := nodeHandler.treeStorage.CheckPermission(ctx, uint(requestData.TreeID), user.User.ID)
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

	storage := nodeHandler.storage

	node, err := storage.CreateNode(ctx, &requestData)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, node)
}

func (nodeHandler *NodeHandler) UpdateAvatar(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	vars := mux.Vars(request)
	nodeID, _ := strconv.Atoi(vars["id"])

	err := request.ParseMultipartForm(2 << 20)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	var preview *multipart.FileHeader

	previewField := request.MultipartForm.File["preview"]
	if len(previewField) != 0 {
		preview = previewField[0]
	} else {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	storage := nodeHandler.storage

	previewPath, err := storage.UpdatePreview(ctx, preview, uint(nodeID))
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	response := &models.UpdatePreviewResponse{
		ID:          uint(nodeID),
		PreviewPath: previewPath,
	}
	responses.SendOkResponse(writer, response)
}

func (nodeHandler *NodeHandler) EditNode(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	vars := mux.Vars(request)
	nodeID, _ := strconv.Atoi(vars["id"])

	var requestData models.EditNodeRequest

	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	storage := nodeHandler.storage

	node, err := storage.EditNode(ctx, &requestData, uint(nodeID))
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, node)
}

func (nodeHandler *NodeHandler) DeleteNode(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	vars := mux.Vars(request)
	nodeID, _ := strconv.Atoi(vars["id"])

	storage := nodeHandler.storage

	err := storage.DeleteNode(ctx, uint(nodeID))
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, &models.SuccessResponse{Success: true})
}
