package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"github.com/yishanzhilu/everest/pkg/bootstrap"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/crypto"
	"github.com/yishanzhilu/everest/pkg/http/server"
)

func main() {
	bootstrap.Boot()
	defer bootstrap.Cleanup()
	guard := crypto.NewJWTGuard(viper.GetString("jwt.secret"), viper.GetDuration("jwt.exp"))
	server := server.NewHTTPServer(
		viper.GetString("runmode"),
		viper.GetString("port"),
		guard,
	)
	server.Start()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	common.Logger.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		common.Logger.Fatal("Server Shutdown: ", err)
	}
	common.Logger.Println("Server exiting")
}
