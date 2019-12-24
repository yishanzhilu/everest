package bootstrap

import (
	"github.com/spf13/viper"
	"github.com/yishanzhilu/api-template/pkg/common"
)

// InitConfig 会使用　viper 读取配置文件
func initConfig(name string) {
	viper.SetConfigName(name)                    // name of config file (without extension)
	viper.AddConfigPath("./configs/")            // for dev, from bin/
	viper.AddConfigPath("../../configs/")        // for test, from pkg/bootstrap/
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		common.Logger.Error("配置初始化失败")
		panic(err)
	}
	common.Logger.Info("配置初始化成功")
}
