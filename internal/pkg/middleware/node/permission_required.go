package node

import (
	authusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	nodeusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/node/usecases"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

func PermissionRequiredMiddleware(storage nodeusecases.NodeStorageInterface,
	authStorage authusecases.AuthStorageInterface) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

			vars := mux.Vars(request)

			nodeID, err := strconv.Atoi(vars["id"])
			if err != nil {
				log.Println(err, responses.StatusInternalServerError)
				responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

				return
			}

			session, _ := request.Cookie("session_id")
			user, _ := authStorage.GetUserBySessionID(ctx, session.Value)

			hasPermission, err := storage.CheckPermission(ctx, uint(nodeID), user.User.ID)
			if err != nil {
				log.Println(err, responses.StatusInternalServerError)
				responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

				return
			}

			if !hasPermission {
				log.Println(responses.ErrForbidden, responses.StatusForbidden)
				responses.SendErrResponse(writer, logger, responses.StatusForbidden, responses.ErrForbidden)

				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}
