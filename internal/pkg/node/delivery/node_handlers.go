package delivery

import (
	authusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/node/usecases"
	"net/http"
)

type NodeHandler struct {
	storage     usecases.NodeStorageInterface
	authStorage authusecases.AuthStorageInterface
}

func NewNodeHandler(storage usecases.NodeStorageInterface, authStorage authusecases.AuthStorageInterface) *NodeHandler {
	return &NodeHandler{
		storage:     storage,
		authStorage: authStorage,
	}
}

func (nodeHandler *NodeHandler) GetDescription(writer http.ResponseWriter, request *http.Request) {}

func (nodeHandler *NodeHandler) CreateNode(writer http.ResponseWriter, request *http.Request) {}

func (nodeHandler *NodeHandler) EditNode(writer http.ResponseWriter, request *http.Request) {}

func (nodeHandler *NodeHandler) DeleteNode(writer http.ResponseWriter, request *http.Request) {}
