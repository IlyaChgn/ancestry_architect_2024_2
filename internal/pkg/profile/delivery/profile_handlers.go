package delivery

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	authusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/usecases"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	storage     usecases.ProfileStorageInterface
	authStorage authusecases.AuthStorageInterface
}

func NewProfileHandler(storage usecases.ProfileStorageInterface,
	authStorage authusecases.AuthStorageInterface) *ProfileHandler {
	return &ProfileHandler{
		storage:     storage,
		authStorage: authStorage,
	}
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

func (profileHandler *ProfileHandler) UpdateProfile(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := utils.GetLoggerFromContext(ctx).With(zap.String("handler", utils.GetFunctionName()))

	authStorage := profileHandler.authStorage
	session, _ := request.Cookie("session_id")

	user, err := authStorage.GetUserBySessionID(ctx, session.Value)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	err = request.ParseMultipartForm(2 << 20)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	var birthdate time.Time

	if request.FormValue("birthdate") != "" {
		birthdate, err = time.Parse("02/01/2006", request.FormValue("birthdate"))
		if err != nil {
			log.Println(err, responses.StatusInternalServerError)
			responses.SendErrResponse(writer, logger, responses.StatusInternalServerError, responses.ErrInternalServer)

			return
		}
	}

	requestData := models.UpdateProfileRequest{
		Email:     request.FormValue("email"),
		Name:      request.FormValue("name"),
		Surname:   request.FormValue("surname"),
		Birthdate: birthdate,
		Gender:    request.FormValue("gender"),
	}

	avatar := request.MultipartForm.File["avatar"]
	if len(avatar) != 0 {
		requestData.Avatar = avatar[0]
	}

	response := models.UpdateProfileResponse{}

	if requestData.Email != "" {
		newUser, errs := authStorage.UpdateEmail(ctx, requestData.Email, user.User.ID)
		if errs != nil {
			log.Println(errs, responses.StatusBadRequest)
			responses.SendSeveralErrsResponse(writer, logger, responses.StatusBadRequest, errs)

			return
		}

		response.User = *newUser
	}

	storage := profileHandler.storage

	profile, err := storage.UpdateProfile(ctx, &requestData, user.User.ID)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

		return
	}

	response.Profile = *profile

	responses.SendOkResponse(writer, response)
}
