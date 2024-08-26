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

func LoginRequiredMiddleware(storage usecases.AuthStorageInterface) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			logger := utils.GetLoggerFromContext(ctx).With(zap.String("middleware", utils.GetFunctionName()))

			session, _ := request.Cookie("session_id")
			if session == nil {
				log.Println("User not authorized", responses.StatusUnauthorized)
				responses.SendErrResponse(writer, logger, responses.StatusUnauthorized, responses.ErrNotAuthorized)

				return
			}

			_, err := storage.GetUserBySessionID(ctx, session.Value)
			if err != nil {
				log.Println("User not authorized", responses.StatusBadRequest)
				responses.SendErrResponse(writer, logger, responses.StatusUnauthorized, responses.ErrNotAuthorized)

				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}
