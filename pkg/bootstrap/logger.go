package bootstrap

import (
	"os"

	"github.com/spf13/viper"
	"github.com/yishanzhilu/api-template/pkg/common"

	"github.com/sirupsen/logrus"
)

func initLogger() {
	if viper.GetString("runmode") != "debug" {
		common.Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		common.Logger.SetFormatter(&logrus.TextFormatter{})
		common.Logger.SetOutput(os.Stdout)
		common.Logger.SetLevel(logrus.DebugLevel)
	}
}
