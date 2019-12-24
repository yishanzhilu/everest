package bootstrap

import (
	"fmt"
	"regexp"

	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/api-template/pkg/common"

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
		common.Logger.WithField("error", err).Fatal("Failed to connect MySQL DB")
	}
	common.Logger.Info("Connect to MySQL DB success")
	common.MySQLClient = db
}
