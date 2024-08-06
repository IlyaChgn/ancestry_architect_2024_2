package recover

import (
	"log"
	"net/http"

	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Panic occurred:", err)
				http.Error(writer, responses.InternalServerError, responses.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
