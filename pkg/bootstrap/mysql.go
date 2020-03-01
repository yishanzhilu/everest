package bootstrap

import (
	"fmt"
	"regexp"

	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/common"

	// init mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func mustConnectMySQL(url, username, password, databasename, parameter string) {
	if url == "" {
		url = fmt.Sprintf("%s:%s@%s?%s", username, password, databasename, parameter)
	}
	var re = regexp.MustCompile(":(.*)@")
	maskedURL := re.ReplaceAllString(url, ":***@")
	common.Logger.WithField("url", maskedURL).Info("Connecting to MySQL DB")

	db, err := gorm.Open("mysql", url)
	if err != nil {
		common.Logger.Error("Failed to connect MySQL DB")
		panic(err)
	}
	common.Logger.Info("Connect to MySQL DB success")

	if viper.GetString("runmode") == "debug" {
		db.LogMode(true)
	}

	common.MySQLClient = db

}
