package auth

import (
	"log"
	"net/http"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func LogoutRequiredMiddleware(storage usecases.AuthStorageInterface) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			logger := utils.GetLoggerFromContext(ctx).With(zap.String("middleware", utils.GetFunctionName()))

			oldSession, _ := request.Cookie("session_id")
			if oldSession != nil {
				if oldUser, _ := storage.GetUserBySessionID(ctx, oldSession.Value); oldUser != nil {
					log.Println("User already authorized", responses.StatusBadRequest)
					responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrAuthorized)

					return
				}
			}

			next.ServeHTTP(writer, request)
		})
	}
}
