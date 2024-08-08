package responses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
	"go.uber.org/zap"
)

const (
	StatusOk = 200

	StatusBadRequest   = 400
	StatusUnauthorized = 401
	StatusForbidden    = 403

	StatusInternalServerError = 500
)

const (
	ErrBadRequest   = "Bad request"
	ErrUnauthorized = "User not authorized"
	ErrAuthorized   = "User already authorized"
	ErrForbidden    = "User have no access to this content"

	ErrInternalServer = "Server error"
)

func newOkResponse(body any) *models.OkResponse {
	return &models.OkResponse{
		Code: StatusOk,
		Body: body,
	}
}

func newErrResponse(code int, status string) *models.ErrResponse {
	return &models.ErrResponse{
		Code:   code,
		Status: status,
	}
}

func sendResponse(writer http.ResponseWriter, response any) {
	serverResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Something went wrong while marshalling JSON", err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}

	_, err = writer.Write(serverResponse)
	if err != nil {
		log.Println("Something went wrong while senddng response", err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}
}

func SendOkResponse(writer http.ResponseWriter, logger *zap.SugaredLogger, body any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(StatusOk)

	response := newOkResponse(body)

	utils.LogHandlerInfo(logger, "success", StatusOk)
	sendResponse(writer, response)
}

func SendErrResponse(writer http.ResponseWriter, logger *zap.SugaredLogger, code int, status string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	response := newErrResponse(code, status)

	utils.LogHandlerError(logger, status, code)
	sendResponse(writer, response)
}
