package grpc

import (
	"context"
	proto "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (manager *AdminManager) GetUsersList(ctx context.Context, empty *emptypb.Empty) (*proto.UserDataList, error) {
	storage := manager.storage

	list, err := storage.GetUsersList(ctx)
	if err != nil {
		return nil, err
	}

	var returningList proto.UserDataList

	for _, val := range *list {
		returningList.Users = append(returningList.Users, &proto.UserData{
			ID:           uint32(val.ID),
			Email:        val.Email,
			PasswordHash: val.PasswordHash,
		})
	}

	return &returningList, nil
}

func (manager *AdminManager) GetAdminBySessionID(ctx context.Context,
	in *proto.SessionRequest) (*proto.UserData, error) {
	storage := manager.storage

	user, err := storage.GetAdminBySessionID(ctx, in.GetSessionID())
	if err != nil {
		return nil, err
	}

	return &proto.UserData{
		ID:           uint32(user.ID),
		Email:        user.Email,
		PasswordHash: "",
	}, nil
}

func (manager *AdminManager) CreateUser(ctx context.Context, in *proto.LoginUserRequest) (*proto.UserData, error) {
	storage := manager.storage

	user, err := storage.CreateUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, err
	}

	return &proto.UserData{
		ID:           uint32(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (manager *AdminManager) DeleteUser(ctx context.Context,
	in *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	storage := manager.storage

	err := storage.DeleteUser(ctx, uint(in.GetID()))
	if err != nil {
		return nil, err
	}

	return &proto.DeleteUserResponse{Success: true}, nil
}
