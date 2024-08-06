package log

import (
	"context"

	"net/http"

	"github.com/google/uuid"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/config"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func CreateLogMiddleware(logger *zap.SugaredLogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), config.LoggerContextKey,
				logger.With(zap.String(string(config.RequestUUIDContextKey), uuid.New().String())))
			request = request.WithContext(ctx)
			next.ServeHTTP(writer, request)
		})
	}
}
