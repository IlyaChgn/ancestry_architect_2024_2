package delivery

import (
	"log"
	"net/http"
	"strconv"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/usecases"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	storage usecases.ProfileStorageInterface
}

func NewProfileHandler(storage usecases.ProfileStorageInterface) *ProfileHandler {
	return &ProfileHandler{
		storage: storage,
	}
}

func (profileHandler *ProfileHandler) UpdateProfile(writer http.ResponseWriter, request *http.Request) {
}

func (profileHandler *ProfileHandler) GetProfile(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	storage := profileHandler.storage

	profile, err := storage.GetProfileByID(ctx, uint(id))
	if err != nil {
		log.Println(responses.ErrBadRequest, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	responses.SendOkResponse(writer, profile)
}
