package routers

import (
	authdel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/delivery"
	authusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	myauth "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/auth"
	mylogger "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/log"
	myrecovery "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/recover"
	nodedel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/node/delivery"
	nodeusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/node/usecases"
	profiledel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/delivery"
	profileusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/usecases"
	treedel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/tree/delivery"
	treeusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/tree/usecases"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.SugaredLogger,
	authStorage authusecases.AuthStorageInterface, profileStorage profileusecases.ProfileStorageInterface,
	treeStorage treeusecases.TreeStorageInterface, nodeStorage nodeusecases.NodeStorageInterface) *mux.Router {
	router := mux.NewRouter()

	router.Use(myrecovery.RecoveryMiddleware)
	router.Use(mylogger.CreateLogMiddleware(logger))

	authHandler := authdel.NewAuthHandler(authStorage, profileStorage)
	profileHandler := profiledel.NewProfileHandler(profileStorage, authStorage)
	treeHandler := treedel.NewTreeHandler(treeStorage, authStorage)
	nodeHandler := nodedel.NewNodeHandler(nodeStorage, authStorage)

	loginRequiredMiddleware := myauth.LoginRequiredMiddleware(authStorage)
	logoutRequiredMiddleware := myauth.LogoutRequiredMiddleware(authStorage)

	rootRouter := router.PathPrefix("/api").Subrouter()
	ServeAuthRouter(rootRouter, authHandler, loginRequiredMiddleware, logoutRequiredMiddleware)
	ServeProfileRouter(rootRouter, profileHandler, loginRequiredMiddleware)
	ServeTreeRouter(rootRouter, treeHandler, loginRequiredMiddleware)
	ServeNodeRouter(rootRouter, nodeHandler, loginRequiredMiddleware)

	return router
}
