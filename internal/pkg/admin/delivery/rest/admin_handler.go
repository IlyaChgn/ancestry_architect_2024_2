package delivery

import grpc "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"

type AdminHandler struct {
	adminClient grpc.AdminClient
}

func NewAdminHandler(adminClient grpc.AdminClient) *AdminHandler {
	return &AdminHandler{
		adminClient: adminClient,
	}
}
