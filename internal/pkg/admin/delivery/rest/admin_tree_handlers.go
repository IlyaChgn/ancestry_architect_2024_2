package delivery

import (
	"encoding/json"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	proto "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

func (adminHandler *AdminHandler) GetTreesList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	userID, err := strconv.Atoi(request.URL.Query().Get("user_id"))
	if err != nil {
		userID = 0
	}

	client := adminHandler.adminClient

	list, err := client.GetTreesList(ctx, &proto.GetTreesListRequest{UserID: uint32(userID)})
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	var trees []models.TreeResponse

	for _, tree := range list.Trees {
		trees = append(trees, models.TreeResponse{
			ID:     uint(tree.ID),
			UserID: uint(tree.UserID),
			Name:   tree.Name,
		})
	}

	responses.SendOkResponse(writer, trees)
}

func (adminHandler *AdminHandler) EditTreeName(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	var requestData models.EditTreeRequest

	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	client := adminHandler.adminClient

	tree, err := client.EditTreeName(ctx, &proto.EditTreeNameRequest{
		TreeID: uint32(requestData.TreeID),
		Name:   requestData.Name,
	})
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, &models.TreeResponse{
		ID:     uint(tree.ID),
		UserID: uint(tree.UserID),
		Name:   tree.Name,
	})
}

func (adminHandler *AdminHandler) GetNodesList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	vars := mux.Vars(request)

	treeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	client := adminHandler.adminClient

	list, err := client.GetNodesList(ctx, &proto.GetNodesListRequest{TreeID: uint32(treeID)})
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	var nodes []models.NodeForAdmin

	for _, node := range list.Nodes {
		nodes = append(nodes, models.NodeForAdmin{
			ID:          uint(node.ID),
			Name:        node.Name,
			Birthdate:   node.Birthdate.AsTime(),
			Deathdate:   node.Deathdate.AsTime(),
			Gender:      node.Gender,
			PreviewPath: node.PreviewPath,
			LayerID:     uint(node.LayerID),
			LayerNum:    int(node.LayerNum),
			TreeID:      uint(node.TreeID),
			UserID:      uint(node.UserID),
			IsDeleted:   node.IsDeleted,
		})
	}

	responses.SendOkResponse(writer, nodes)
}
