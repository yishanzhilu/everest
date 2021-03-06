package bootstrap

import (
	"os"

	"github.com/spf13/viper"
	"github.com/yishanzhilu/everest/pkg/common"

	"github.com/sirupsen/logrus"
)

func initLogger() {
	if viper.GetString("runmode") != "debug" {
		common.Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		common.Logger.SetFormatter(&logrus.TextFormatter{})
		common.Logger.SetOutput(os.Stdout)
		common.Logger.SetReportCaller(true)
		common.Logger.SetLevel(logrus.DebugLevel)
	}
}
