package delivery

import (
	grpc "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"
	mysession "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/session"
	"net/http"
	"time"
)

type AdminHandler struct {
	adminClient grpc.AdminClient
}

func NewAdminHandler(adminClient grpc.AdminClient) *AdminHandler {
	return &AdminHandler{
		adminClient: adminClient,
	}
}

func (adminHandler *AdminHandler) createSession(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     "admin_sid",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(mysession.AdminSessionDuration),
		HttpOnly: true,
		SameSite: 4,
		Secure:   true,
	}
}
