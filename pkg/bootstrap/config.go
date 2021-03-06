package bootstrap

import (
	"github.com/spf13/viper"
	"github.com/yishanzhilu/everest/pkg/common"
)

// InitConfig 会使用　viper 读取配置文件, name 为 /configs/ 下的文件名
// 如果存在 name.local,
func initConfig(name string) {
	viper.SetConfigName(name)                    // name of config file (without extension)
	viper.AddConfigPath("./configs/")            // for dev, from bin/
	viper.AddConfigPath("../../configs/")        // for test, from pkg/bootstrap/
	viper.AddConfigPath("/etc/everest/")         // for prod
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		common.Logger.Error("配置初始化失败")
		panic(err)
	}
	common.Logger.Info("配置初始化成功")
	common.JWTSecret = viper.GetString("jwt.secret")
	common.JWTTokenExpDuration = viper.GetDuration("jwt.token-exp-duration")
}
