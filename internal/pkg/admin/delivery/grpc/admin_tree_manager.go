package grpc

import (
	"context"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	proto "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (manager *AdminManager) GetNodesList(ctx context.Context,
	in *proto.GetNodesListRequest) (*[]proto.NodeData, error) {
	storage := manager.storage

	list, err := storage.GetNodesList(ctx, uint(in.GetTreeID()))
	if err != nil {
		return nil, err
	}

	var returningList []proto.NodeData

	for _, val := range *list {
		returningList = append(returningList, proto.NodeData{
			ID:          uint32(val.ID),
			Name:        val.Name,
			Birthdate:   timestamppb.New(*val.Birthdate),
			Deathdate:   timestamppb.New(*val.Deathdate),
			Gender:      val.Gender,
			PreviewPath: val.PreviewPath,
			LayerID:     uint32(val.LayerID),
			LayerNum:    int32(val.LayerNum),
			TreeID:      uint32(val.TreeID),
			UserID:      uint32(val.UserID),
			IsDeleted:   val.IsDeleted,
		})
	}

	return &returningList, nil
}

func (manager *AdminManager) EditTreeName(ctx context.Context,
	in *proto.EditTreeNameRequest) (*proto.TreeData, error) {
	storage := manager.storage

	tree, err := storage.EditTreeName(ctx, uint(in.GetTreeID()), in.GetName())
	if err != nil {
		return nil, err
	}

	return &proto.TreeData{
		ID:     uint32(tree.ID),
		UserID: uint32(tree.UserID),
		Name:   tree.Name,
	}, nil
}

func (manager *AdminManager) GetTreesList(ctx context.Context,
	in *proto.GetTreesListRequest) (*[]proto.TreeData, error) {
	storage := manager.storage

	var (
		list *[]models.TreeResponse
		err  error
	)

	if in.GetUserID() != 0 {
		list, err = storage.GetTreesListByUserID(ctx, uint(in.GetUserID()))
	} else {
		list, err = storage.GetTreesList(ctx)
	}

	if err != nil {
		return nil, err
	}

	var returningList []proto.TreeData

	for _, val := range *list {
		returningList = append(returningList, proto.TreeData{
			ID:     uint32(val.ID),
			UserID: uint32(val.UserID),
			Name:   val.Name,
		})
	}

	return &returningList, nil
}
