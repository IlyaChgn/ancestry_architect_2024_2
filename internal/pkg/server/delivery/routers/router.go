package routers

import (
	grpc "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/grpc/protobuf"
	admindel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/admin/delivery/rest"
	authdel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/delivery"
	authusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/auth/usecases"
	myauth "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/auth"
	mylogger "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/log"
	mynode "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/node"
	myrecovery "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/middleware/recover"
	nodedel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/node/delivery"
	nodeusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/node/usecases"
	profiledel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/delivery"
	profileusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/profile/usecases"
	routers "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery/routers/admin"
	approuters "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/delivery/routers/app"
	treedel "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/tree/delivery"
	treeusecases "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/tree/usecases"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.SugaredLogger,
	authStorage authusecases.AuthStorageInterface, profileStorage profileusecases.ProfileStorageInterface,
	treeStorage treeusecases.TreeStorageInterface, nodeStorage nodeusecases.NodeStorageInterface,
	adminClient grpc.AdminClient) *mux.Router {
	router := mux.NewRouter()

	router.Use(myrecovery.RecoveryMiddleware)
	router.Use(mylogger.CreateLogMiddleware(logger))

	authHandler := authdel.NewAuthHandler(authStorage, profileStorage)
	profileHandler := profiledel.NewProfileHandler(profileStorage, authStorage)
	treeHandler := treedel.NewTreeHandler(treeStorage, authStorage)
	nodeHandler := nodedel.NewNodeHandler(nodeStorage, treeStorage, authStorage)
	adminHandler := admindel.NewAdminHandler(adminClient)

	loginRequiredMiddleware := myauth.LoginRequiredMiddleware(authStorage)
	logoutRequiredMiddleware := myauth.LogoutRequiredMiddleware(authStorage)
	permissionRequiredMiddleware := mynode.PermissionRequiredMiddleware(nodeStorage, authStorage)

	apiRouter := router.PathPrefix("/api").Subrouter()
	approuters.ServeAppAuthRouter(apiRouter, authHandler, loginRequiredMiddleware, logoutRequiredMiddleware)
	approuters.ServeAppProfileRouter(apiRouter, profileHandler, loginRequiredMiddleware)
	approuters.ServeAppTreeRouter(apiRouter, treeHandler, loginRequiredMiddleware)
	approuters.ServeAppNodeRouter(apiRouter, nodeHandler, loginRequiredMiddleware, permissionRequiredMiddleware)

	adminRouter := router.PathPrefix("/admin").Subrouter()
	routers.ServeAdminAuthRouter(adminRouter, adminHandler)
	routers.ServeAdminTreeRouter(adminRouter, adminHandler)

	return router
}
