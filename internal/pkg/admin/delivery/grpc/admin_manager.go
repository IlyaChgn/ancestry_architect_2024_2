package grpc

import (
	protobuf "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/usecases"
)

type AdminManager struct {
	protobuf.UnimplementedAdminServer
	storage usecases.AdminStorageInterface
}

func NewAdminManager(adminStorage usecases.AdminStorageInterface) *AdminManager {
	return &AdminManager{
		storage: adminStorage,
	}
}
