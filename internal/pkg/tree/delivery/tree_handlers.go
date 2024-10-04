package delivery

import (
	authusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/tree/usecases"
	"net/http"
)

type TreeHandler struct {
	storage     usecases.TreeStorageInterface
	authStorage authusecases.AuthStorageInterface
}

func NewTreeHandler(storage usecases.TreeStorageInterface, authStorage authusecases.AuthStorageInterface) *TreeHandler {
	return &TreeHandler{
		storage:     storage,
		authStorage: authStorage,
	}
}

func (treeHandler *TreeHandler) GetTree(writer http.ResponseWriter, request *http.Request) {}

func (treeHandler *TreeHandler) GetCreatedTreesList(writer http.ResponseWriter, request *http.Request) {
}

func (treeHandler *TreeHandler) GetAvailableTreesList(writer http.ResponseWriter, request *http.Request) {
}

func (treeHandler *TreeHandler) CreateTree(writer http.ResponseWriter, request *http.Request) {}

func (treeHandler *TreeHandler) AddPermission(writer http.ResponseWriter, request *http.Request) {}
