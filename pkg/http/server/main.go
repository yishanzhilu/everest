package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/crypto"
	"github.com/yishanzhilu/everest/pkg/http/server/middleware"
	"github.com/yishanzhilu/everest/pkg/repo/mysql"
	"github.com/yishanzhilu/everest/pkg/user"
	"gopkg.in/resty.v1"
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

	// workspaceRouter := v1.Group("workspace")
	// bootWorkspace(workspaceRouter)

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

func bootUser(userRouter *gin.RouterGroup, guard *crypto.JWTGuard) {
	userHandler := constructNewUserHandler(common.Logger, common.MySQLClient, guard)
	userHandler.RegisterPublicRoutes(userRouter)
	userRouter.Use(middleware.Authenticate())
	userHandler.RegisterPrivateRoutes(userRouter)
}

// constructNewUserHandler ..
func constructNewUserHandler(logger *logrus.Logger, db *gorm.DB, guard *crypto.JWTGuard) user.Handler {
	githubClient := resty.New().
		SetTimeout(30*time.Second).
		SetHostURL("https://github.com").
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"client_id":     viper.GetString("github.client.id"),
			"client_secret": viper.GetString("github.client.secret"),
		})
	gr := user.NewGithubRepo(githubClient, logger)
	mr := mysql.NewMysqlUserRepository(db)
	userService := user.NewUserService(mr, gr)
	userHandler := user.NewHandler(userService, guard)
	return userHandler
}
