package bootstrap

import (
	"sync"

	"github.com/spf13/viper"
	"github.com/yishanzhilu/api-template/pkg/common"
	"github.com/yishanzhilu/api-template/pkg/http/server"
)

// Boot will bootstarp the program, commonly used with Close
//
// func main() {
// 	bootstrap.Boot()
// 	defer bootstrap.Close()
// }
func Boot() {

	initConfig("viper")
	initLogger()
	initHTTPClient()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		mustConnectRedis(
			viper.GetString("redis.url"),
			viper.GetString("redis.password"),
		)
		wg.Done()
	}()
	go func() {
		mustConnectMySQL(
			viper.GetString("mysql.url"),
			viper.GetString("mysql.username"),
			viper.GetString("mysql.password"),
			viper.GetString("mysql.databasename"),
			viper.GetString("mysql.parameter"),
		)
		wg.Done()
	}()
	wg.Wait()
	server.Start()
}

// Close is used for cleaning up
func Close() {
	common.MySQLClient.Close()
	common.RedisClient.Close()
}
