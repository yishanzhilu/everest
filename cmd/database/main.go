package main

import (
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/bootstrap"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/models"
)

func main() {
	bootstrap.Boot()
	defer bootstrap.Cleanup()
	common.MySQLClient.LogMode(true)
	common.MySQLClient.AutoMigrate(&models.UserModel{})
	common.MySQLClient.AutoMigrate(&models.GoalModel{})
	common.MySQLClient.AutoMigrate(&models.MissionModel{})
	common.MySQLClient.AutoMigrate(&models.RecordModel{})
	common.MySQLClient.AutoMigrate(&models.TodoModel{})
	LoadCSV(common.MySQLClient, "record_models", true)
}

// LoadCSV 会把 /tmp/seeds/:name 路径中的csv加载到表表中，需要保证csv字段顺序和model一致
func LoadCSV(db *gorm.DB, name string, updateTimestamp bool) {
	csvPath := "./cmd/database/seeds/" + name + ".csv"
	mysql.RegisterLocalFile(csvPath)
	sql := "LOAD DATA LOCAL INFILE '" + csvPath + "' INTO TABLE " + name + " fields terminated BY \",\"lines terminated BY \"\n\" IGNORE 1 LINES"
	if updateTimestamp {
		sql += " set deleted_at = NULL"
	}
	err := db.Exec(sql).Error
	if err != nil {
		log.Fatal("LOAD DATA LOCAL INFILE ", err)
	}
}
