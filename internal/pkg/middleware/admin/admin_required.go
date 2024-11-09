package admin

import (
	grpc "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func AdminRequiredMiddleware(client grpc.AdminClient) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			logger := utils.GetLoggerFromContext(ctx).With(zap.String("middleware", utils.GetFunctionName()))

			session, _ := request.Cookie("admin_sid")
			if session == nil {
				log.Println("User not authorized", responses.StatusUnauthorized)
				responses.SendErrResponse(writer, logger, responses.StatusUnauthorized, responses.ErrNotAuthorized)

				return
			}

			_, err := client.GetAdminBySessionID(ctx, &grpc.SessionRequest{SessionID: session.Value})
			if err != nil {
				log.Println(err, responses.StatusForbidden)
				responses.SendErrResponse(writer, logger, responses.StatusForbidden, responses.ErrForbidden)

				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}
