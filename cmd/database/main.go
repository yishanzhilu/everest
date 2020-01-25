package main

import (
	"github.com/yishanzhilu/everest/pkg/bootstrap"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/models"
)

func main() {
	bootstrap.Boot()
	defer bootstrap.Cleanup()
	common.MySQLClient.AutoMigrate(&models.UserModel{})
	common.MySQLClient.AutoMigrate(&models.GoalModel{})
	common.MySQLClient.AutoMigrate(&models.MissionModel{})
	common.MySQLClient.AutoMigrate(&models.TaskModel{})
}
