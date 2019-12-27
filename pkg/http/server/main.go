package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/crypto"
	"github.com/yishanzhilu/everest/pkg/http/server/middleware"
	"github.com/yishanzhilu/everest/pkg/storage/mysql"
	"github.com/yishanzhilu/everest/pkg/user"
	"github.com/yishanzhilu/everest/pkg/workspace"
)

// Server is a program that can start server http request at a specifc port.
type Server interface {
	Start()
	Shutdown(context.Context) error
}

// HTTPServer .
type HTTPServer struct {
	runmode string
	port    string
	guard   *crypto.JWTGuard
	server  *http.Server
}

// NewHTTPServer NewHttpServer.
func NewHTTPServer(runmode, port string, guard *crypto.JWTGuard) Server {
	return &HTTPServer{
		runmode,
		port,
		guard,
		nil,
	}
}

// Start implementation.
func (s *HTTPServer) Start() {
	gin.SetMode(s.runmode)
	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.GinLogger())
	router.Use(gin.Recovery())

	healthCheck(router)

	v1 := router.Group("/api/v1")
	v1.Use(middleware.AssignGuard(s.guard))

	workspaceRouter := v1.Group("workspace")
	bootWorkspace(workspaceRouter)

	userRouter := v1.Group("user")
	bootUser(userRouter, s.guard)

	s.server = &http.Server{
		Addr:    s.port,
		Handler: router,
	}
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		common.Logger.WithField("err", err).Fatal("Server Shutdown")
	}
}

// Shutdown implementation.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func bootWorkspace(workspaceRouter *gin.RouterGroup) {
	workprofileMysqlRepo := mysql.NewMysqlWorkprofileRepository(common.MySQLClient)
	workprofileService := workspace.NewWorkprofileService(workprofileMysqlRepo)
	workprofileHandler := workspace.NewWorkprofileHandler(workprofileService)
	workprofileHandler.RegisterPublicRoutes(workspaceRouter)
	workspaceRouter.Use(middleware.Authenticate())
	workprofileHandler.RegisterPrivateRoutes(workspaceRouter)
}

func bootUser(userRouter *gin.RouterGroup, guard *crypto.JWTGuard) {
	userHandler := user.ConstructNewUserHandler(common.Logger, common.MySQLClient, guard)
	userHandler.RegisterPublicRoutes(userRouter)
	userRouter.Use(middleware.Authenticate())
	userHandler.RegisterPrivateRoutes(userRouter)
}
