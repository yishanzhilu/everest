package bootstrap

import (
	"os"
	"sync"

	"github.com/spf13/viper"
	"github.com/yishanzhilu/everest/pkg/common"
)

// Boot will bootstarp the program, commonly used with Cleanup
//
// func main() {
// 	bootstrap.Boot()
// 	defer bootstrap.Cleanup()
// }
func Boot() {
	filename := os.Getenv("EVEREST_CONFIG_FILE_NAME")
	if filename == "" {
		filename = "viper.local"
	}
	initConfig(filename)
	initLogger()
	initHTTPClient()
	var wg sync.WaitGroup
	wg.Add(1)
	// go func() {
	// 	mustConnectRedis(
	// 		viper.GetString("redis.url"),
	// 		viper.GetString("redis.password"),
	// 	)
	// 	wg.Done()
	// }()
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
}

// Cleanup is used for cleaning up
func Cleanup() {
	common.MySQLClient.Close()
	// common.RedisClient.Close()
}
