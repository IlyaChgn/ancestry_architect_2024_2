package responses_test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/models"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/config"
	responses "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery"
	logger "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/usecases"
)

func TestOkResponse(t *testing.T) {
	writer := httptest.NewRecorder()
	data := "test case"

	logger, err := logger.NewLogger(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "))
	if err != nil {
		t.Fatalf("Error while creating logger, %v", err)
	}

	responses.SendOkResponse(writer, logger, data)

	if code := writer.Code; code != responses.StatusOk {
		t.Errorf("Expected status 200, got: %d", code)
	}

	res := writer.Result()
	defer res.Body.Close()

	var resultBody models.OkResponse
	err = json.NewDecoder(res.Body).Decode(&resultBody)
	if err != nil {
		t.Fatalf("Expected error to be nil got %v", err)
	}

	if resultBody.Body != "test case" {
		t.Errorf("Expected 'test case' got %v", resultBody.Body)
	}
}

func TestErrResponse(t *testing.T) {
	writer := httptest.NewRecorder()

	logger, err := logger.NewLogger(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "))
	if err != nil {
		t.Fatalf("Error while creating logger, %v", err)
	}

	responses.SendErrResponse(writer, logger, responses.StatusBadRequest, responses.ErrBadRequest)

	if code := writer.Code; code != responses.StatusBadRequest {
		t.Errorf("Expected status 400, got: %d", code)
	}

	res := writer.Result()
	defer res.Body.Close()

	var resultBody models.ErrResponse
	err = json.NewDecoder(res.Body).Decode(&resultBody)
	if err != nil {
		t.Fatalf("Expected error to be nil got %v", err)
	}

	if resultBody.Status != responses.ErrBadRequest {
		t.Errorf("Expected 'Bad request' got %v", resultBody.Status)
	}
}
