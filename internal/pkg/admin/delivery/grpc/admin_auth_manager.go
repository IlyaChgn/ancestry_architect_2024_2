package grpc

import (
	"context"
	proto "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/google/uuid"
)

func (manager *AdminManager) Login(ctx context.Context, in *proto.LoginUserRequest) (*proto.UserAuthResponse, error) {
	storage := manager.storage

	admin, err := storage.GetAdminByEmail(ctx, in.GetEmail())
	if err != nil {
		return nil, err
	}

	if !utils.CheckPassword(in.GetPassword(), admin.PasswordHash) {
		return nil, err
	}

	sessionID := uuid.NewString()

	err = storage.CreateSession(ctx, sessionID, admin.ID)
	if err != nil {
		return nil, err
	}

	return &proto.UserAuthResponse{
		ID:           uint32(admin.ID),
		Email:        admin.Email,
		PasswordHash: "",
		SessionID:    sessionID,
	}, nil
}

func (manager *AdminManager) Logout(ctx context.Context, in *proto.SessionRequest) (*proto.UserAuthResponse, error) {
	storage := manager.storage

	err := storage.RemoveSession(ctx, in.GetSessionID())
	if err != nil {
		return nil, err
	}

	return &proto.UserAuthResponse{}, nil
}

func (manager *AdminManager) EditPassword(ctx context.Context, in *proto.EditPasswordRequest) (*proto.UserData, error) {
	storage := manager.storage

	user, err := storage.UpdatePassword(ctx, uint(in.GetID()), in.GetPassword())
	if err != nil {
		return nil, err
	}

	return &proto.UserData{
		ID:           uint32(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (manager *AdminManager) GetUsersList(ctx context.Context) (*[]proto.UserData, error) {
	storage := manager.storage

	list, err := storage.GetUsersList(ctx)
	if err != nil {
		return nil, err
	}

	var returningList []proto.UserData

	for _, val := range *list {
		returningList = append(returningList, proto.UserData{
			ID:           uint32(val.ID),
			Email:        val.Email,
			PasswordHash: val.PasswordHash,
		})
	}

	return &returningList, nil
}
